package bg

import (
	"time"

	"gophers.dev/pkgs/loggy"
	"gophers.dev/pkgs/repeat/x"

	"oss.indeed.com/go/modprox/pkg/clients/registry"
	"oss.indeed.com/go/modprox/pkg/coordinates"
	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/proxy/internal/modules/get"
	"oss.indeed.com/go/modprox/proxy/internal/modules/store"
	"oss.indeed.com/go/modprox/proxy/internal/problems"
)

type Options struct {
	// Frequency determines how often the worker will check in
	// with the registry, looking for new modules that need to be
	// downloaded by this instance of the proxy. A typical value
	// would be something like 30 seconds - not too slow, but also
	// not spamming the network with polling traffic.
	Frequency time.Duration
}

// A Worker runs in the background, polling the registry for new
// modules that need to be downloaded, and downloading those modules
// as needed.
type Worker interface {
	Start(options Options)
}

type worker struct {
	registryClient    registry.Client
	emitter           stats.Sender
	dlTracker         problems.Tracker
	index             store.Index
	store             store.ZipStore
	downloader        get.Downloader
	registryRequester get.RegistryAPI
	log               loggy.Logger
}

func New(
	emitter stats.Sender,
	dlTracker problems.Tracker,
	index store.Index,
	store store.ZipStore,
	registryRequester get.RegistryAPI,
	downloader get.Downloader,
) Worker {
	return &worker{
		emitter:           emitter,
		dlTracker:         dlTracker,
		index:             index,
		store:             store,
		downloader:        downloader,
		registryRequester: registryRequester,
		log:               loggy.New("bg-worker"),
	}
}

func (w *worker) Start(options Options) {
	go func() {
		_ = x.Interval(options.Frequency, func() error {
			if err := w.loop(); err != nil {
				w.log.Errorf("worker loop iteration had error: %v", err)
				// never return an error, which would stop the worker
				// instead, we remain hopeful the next iteration will work
			}
			return nil
		})
	}()
}

func (w *worker) loop() error {
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

func (w *worker) acquireMods() ([]coordinates.SerialModule, error) {
	ids, err := w.index.IDs()
	if err != nil {
		return nil, err
	}

	mods, err := w.registryRequester.ModulesNeeded(ids)
	if err != nil {
		w.log.Errorf("failed to acquire list of needed mods from registry, %v", err)
		return nil, err
	}
	w.log.Infof("acquired list of %d mods from registry", len(mods))

	for _, mod := range mods {
		w.log.Tracef("- %s @ %s", mod.Source, mod.Version)
	}

	for _, mod := range mods {
		// only download mod if we do not already have it
		exists, indexID, err := w.index.Contains(mod.Module)
		if err != nil {
			w.log.Errorf("problem with index lookups: %v", err)
			continue // may as well try the others
		}

		if exists {
			w.log.Tracef("already have %s, not going to download it again", mod)
			// set indexID to newID if they do not match
			if indexID != mod.SerialID {
				w.log.Infof(
					"indexed ID of %d for %s does not match ID %d, will update",
					indexID,
					mod,
					mod.SerialID,
				)
				if err = w.index.UpdateID(mod); err != nil {
					w.log.Errorf("problem updating index ID: %v", err)
				}
			}
			continue // move on to the next one
		}

		if err := w.downloader.Download(mod); err != nil {
			w.log.Errorf("failed to download %s, %v", mod, err)
			w.dlTracker.Set(problems.Create(mod.Module, err))
			continue // may as well try the others
		}
		w.log.Tracef("downloaded %s!", mod)
	}

	return mods, nil
}
