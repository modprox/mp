package get

import (
	"fmt"
	"time"

	"oss.indeed.com/go/modprox/pkg/clients/zips"
	"oss.indeed.com/go/modprox/pkg/coordinates"
	"oss.indeed.com/go/modprox/pkg/loggy"
	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/pkg/repository"
	"oss.indeed.com/go/modprox/pkg/upstream"
	"oss.indeed.com/go/modprox/proxy/internal/modules/store"
)

//go:generate go run github.com/gojuno/minimock/cmd/minimock -g -i Downloader -s _mock.go

type Downloader interface {
	Download(module coordinates.SerialModule) error
}

func New(
	proxyClient zips.ProxyClient,
	upstreamClient zips.UpstreamClient,
	resolver upstream.Resolver,
	store store.ZipStore,
	index store.Index,
	emitter stats.Sender,
) Downloader {
	return &downloader{
		proxyClient:    proxyClient,
		upstreamClient: upstreamClient,
		resolver:       resolver,
		store:          store,
		index:          index,
		emitter:        emitter,
		log:            loggy.New("downloader"),
	}
}

type downloader struct {
	proxyClient    zips.ProxyClient
	upstreamClient zips.UpstreamClient
	resolver       upstream.Resolver
	store          store.ZipStore
	index          store.Index
	emitter        stats.Sender
	log            loggy.Logger
}

func (d *downloader) downloadFromProxy(mod coordinates.SerialModule) (repository.Blob, error) {
	d.log.Infof("going to download from proxy: %s", mod.String())

	// download the well-formed zip from the proxy
	start := time.Now()
	blob, err := d.proxyClient.Get(mod.Module)
	if err != nil {
		return nil, err
	}

	d.emitter.GaugeMS("download-mod-elapsed-ms", start)
	d.log.Infof("downloaded upstream blob of size: %d", len(blob))

	// no need to re-write, this is already a correctly formatted zip
	return blob, nil
}

func (d *downloader) downloadFromUpstream(mod coordinates.SerialModule) (repository.Blob, error) {
	d.log.Infof("going to download from upstream: %s", mod.String())

	request, err := d.resolver.Resolve(mod.Module)
	if err != nil {
		return nil, err
	}

	// download the raw-zip from the upstream source
	start := time.Now()
	blob, err := d.upstreamClient.Get(request)
	if err != nil {
		return nil, err
	}

	d.emitter.GaugeMS("download-mod-elapsed-ms", start)
	d.log.Infof("downloaded upstream blob of size: %d", len(blob))

	rewritten, err := zips.Rewrite(mod.Module, blob)
	if err != nil {
		d.log.Errorf("failed to rewrite blob for %s, %v", mod, err)
		return nil, err
	}

	return rewritten, nil
}

func (d *downloader) storeBlob(mod coordinates.SerialModule, blob repository.Blob) error {

	if err := d.store.PutZip(mod.Module, blob); err != nil {
		d.log.Errorf("failed to save blob to zip store for %s, %v", mod, err)
		return err
	}

	modFile, exists, err := blob.ModFile()
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

	d.log.Tracef("stored %s, was %d bytes", mod, len(blob))

	return nil
}

func (d *downloader) Download(mod coordinates.SerialModule) error {
	useProxy, err := d.resolver.UseProxy(mod.Module)
	if err != nil {
		d.log.Errorf("could not decide on using proxy:", err)
		return err
	}
	{
		var (
			blob repository.Blob
			err  error
		)
		switch useProxy {
		case true:
			blob, err = d.downloadFromProxy(mod)
		default:
			blob, err = d.downloadFromUpstream(mod)
		}

		if err != nil {
			d.log.Errorf("failed to download %s: %v", mod, err)
			return err
		}

		return d.storeBlob(mod, blob)
	}
}

func emptyModFile(mod coordinates.SerialModule) string {
	return fmt.Sprintf("module %s\n", mod.Module.Source)
}
