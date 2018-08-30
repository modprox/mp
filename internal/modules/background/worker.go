package background

import (
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
	store          store.Store
	resolver       upstream.Resolver
	downloader     zips.Client
	log            loggy.Logger
}

func NewReloader(
	options Options,
	registryClient registry.Client,
	store store.Store,
	resolver upstream.Resolver,
	downloader zips.Client,
) Reloader {
	return &reloadWorker{
		options:        options,
		registryClient: registryClient,
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
	mods, err := w.registryClient.ModInfos()
	if err != nil {
		return nil, err
	}
	w.log.Infof("acquired list of %d mods from registry", len(mods))

	for _, mod := range mods {
		w.log.Tracef("- %s @ %s", mod.Source, mod.Version)
	}

	for _, mod := range mods {
		_ = w.download(mod)
	}

	return mods, nil
}

func (w *reloadWorker) download(mod repository.ModInfo) error {
	request, err := w.resolver.Resolve(mod)
	if err != nil {
		return err
	}

	w.log.Infof("will download %s", request.URI())

	// actually download it
	blob, err := w.downloader.Get(request)
	if err != nil {
		return err
	}

	w.log.Infof("downloaded blob of size: %d", len(blob))

	return nil
}
