package background

import (
	"fmt"
	"time"

	"github.com/modprox/libmodprox/clients/registry"
	"github.com/modprox/libmodprox/clients/zips"
	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/libmodprox/repository"
	"github.com/modprox/libmodprox/upstream"
	"github.com/modprox/modprox-proxy/internal/modules/store"

	"github.com/shoenig/toolkit"
)

type Options struct {
	Frequency time.Duration
}

type Reloader interface {
	Start()
}

type reloadWorker struct {
	options        Options
	registryClient registry.Client
	index          store.Index
	store          store.ZipStore
	resolver       upstream.Resolver
	downloader     zips.Client
	log            loggy.Logger
}

func NewReloader(
	options Options,
	registryClient registry.Client,
	index store.Index,
	store store.ZipStore,
	resolver upstream.Resolver,
	downloader zips.Client,
) Reloader {
	return &reloadWorker{
		options:        options,
		registryClient: registryClient,
		index:          index,
		store:          store,
		resolver:       resolver,
		downloader:     downloader,
		log:            loggy.New("reload-worker"),
	}
}

func (w *reloadWorker) Start() {
	go toolkit.Interval(w.options.Frequency, func() error {
		if err := w.loop(); err != nil {
			w.log.Errorf("worker loop iteration had error: %v", err)
			// never return an error, which would stop the worker
			// instead, we remain hopeful the next iteration will work
		}
		return nil
	})
}

func (w *reloadWorker) loop() error {
	w.log.Infof("worker loop starting")

	mods, err := w.acquireMods()
	if err != nil {
		return err
	}

	_ = mods

	// we have a list of modules already downloaded to fs
	// we have a list of modules from registry that we want
	// do a diff, finding:
	// - modules we have but do not need anymore
	// - modules we need but to not have yet
	// then prune modules we do not want
	// then DL and save modules we do want
	// also, take into account redirects and such
	return nil
}

func (w *reloadWorker) acquireMods() ([]repository.ModInfo, error) {
	mods, err := w.registryClient.Modules()
	if err != nil {
		return nil, err
	}
	w.log.Infof("acquired list of %d mods from registry", len(mods))

	for _, mod := range mods {
		w.log.Tracef("- %s @ %s", mod.Source, mod.Version)
	}

	for _, mod := range mods {
		// only download mod if we do not already have it
		exists, err := w.index.Contains(mod)
		if err != nil {
			w.log.Errorf("problem with index lookups: %v", err)
			continue // may as well try the others
		}

		if exists {
			w.log.Tracef("already have %s, not going to download it again", mod)
			continue // move on to the next one
		}

		if err := w.download(mod); err != nil {
			w.log.Errorf("failed to download %s, %v", mod, err)
			continue // may as well try the others
		}
		w.log.Tracef("downloaded %s!", mod)
	}

	return mods, nil
}

func (w *reloadWorker) download(mod repository.ModInfo) error {
	request, err := w.resolver.Resolve(mod)
	if err != nil {
		return err
	}

	w.log.Infof("going to download %s", request.URI())

	// actually download it
	blob, err := w.downloader.Get(request)
	if err != nil {
		return err
	}

	w.log.Infof("downloaded blob of size: %d", len(blob))

	rewritten, err := zips.Rewrite(mod, blob)
	if err != nil {
		w.log.Errorf("failed to rewrite blob for %s, %v", mod, err)
		return err
	}

	if err := w.store.PutZip(mod, rewritten); err != nil {
		w.log.Errorf("failed to save blob to zip store for %s, %v", mod, err)
		return err
	}

	modFile, exists, err := rewritten.ModFile()
	if err != nil {
		w.log.Errorf("failed to re-read re-written zip file for %s, %v", mod, err)
		return err
	}
	if !exists {
		modFile = emptyModFile(mod)
	}

	ma := store.ModuleAddition{
		Mod:      mod,
		UniqueID: 666, // todo: send from registry
		ModFile:  modFile,
	}

	if err := w.index.Put(ma); err != nil {
		w.log.Errorf("failed to updated index for %s, %v", mod, err)
		return err
	}

	return nil
}

func emptyModFile(mod repository.ModInfo) string {
	return fmt.Sprintf("module %s\n", mod.Source)
}
