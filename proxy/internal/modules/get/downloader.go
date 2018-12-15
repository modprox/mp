package get

import (
	"fmt"
	"time"

	"github.com/modprox/mp/pkg/clients/zips"
	"github.com/modprox/mp/pkg/coordinates"
	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/pkg/metrics/stats"
	"github.com/modprox/mp/pkg/upstream"
	"github.com/modprox/mp/proxy/internal/modules/store"
)

//go:generate mockery3 -interface Downloader -package fetchtest

type Downloader interface {
	Download(module coordinates.SerialModule) error
}

func New(
	zipClient zips.Client,
	resolver upstream.Resolver,
	store store.ZipStore,
	index store.Index,
	emitter stats.Sender,
) Downloader {
	return &downloader{
		zipClient: zipClient,
		resolver:  resolver,
		store:     store,
		index:     index,
		emitter:   emitter,
		log:       loggy.New("downloader"),
	}
}

type downloader struct {
	zipClient zips.Client
	resolver  upstream.Resolver
	store     store.ZipStore
	index     store.Index
	emitter   stats.Sender
	log       loggy.Logger
}

func (d *downloader) Download(mod coordinates.SerialModule) error {
	request, err := d.resolver.Resolve(mod.Module)
	if err != nil {
		return err
	}

	d.log.Infof("going to download %s", request.URI())

	// actually download it
	start := time.Now()
	blob, err := d.zipClient.Get(request) // todo: rc bug in there
	if err != nil {
		return err
	}

	d.emitter.GaugeMS("download-mod-elapsed-ms", start)
	d.log.Infof("downloaded blob of size: %d", len(blob))

	rewritten, err := zips.Rewrite(mod.Module, blob)
	if err != nil {
		d.log.Errorf("failed to rewrite blob for %s, %v", mod, err)
		return err
	}

	if err := d.store.PutZip(mod.Module, rewritten); err != nil {
		d.log.Errorf("failed to save blob to zip store for %s, %v", mod, err)
		return err
	}

	modFile, exists, err := rewritten.ModFile()
	if err != nil {
		d.log.Errorf("failed to re-read re-written zip file for %s, %v", mod, err)
		return err
	}
	if !exists {
		modFile = emptyModFile(mod)
	}

	ma := store.ModuleAddition{
		Mod:      mod.Module,
		UniqueID: mod.SerialID,
		ModFile:  modFile,
	}

	if err := d.index.Put(ma); err != nil {
		d.log.Errorf("failed to updated index for %s, %v", mod, err)
		return err
	}

	return nil
}

func emptyModFile(mod coordinates.SerialModule) string {
	return fmt.Sprintf("module %s\n", mod.Module.Source)
}
