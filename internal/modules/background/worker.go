package background

import (
	"time"

	"github.com/modprox/libmodprox/loggy"
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
	options Options
	store   store.Store
	log     loggy.Logger
}

func NewReloader(options Options, store store.Store) Reloader {
	return &reloadWorker{
		options: options,
		store:   store,
		log:     loggy.New("reload-worker"),
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
	// we have a list of modules already downloaded to fs
	// we have a list of modules from registry that we want
	// do a diff, finding:
	// - modules we have but do not need anymore
	// - modules we need but to not have yet
	// then prune modules we do not want
	// then DL and save modules we do want
	return nil
}
