package background

import (
	"log"
	"time"

	"github.com/shoenig/toolkit"

	"github.com/modprox/modprox-proxy/internal/modules/store"
)

type Options struct {
	Frequency time.Duration
}

type Reloader interface {
	Start()
}

type worker struct {
	options Options
	store   store.Store
}

func NewReloader(options Options, store store.Store) Reloader {
	return &worker{
		options: options,
		store:   store,
	}
}

func (w *worker) Start() {
	go toolkit.Interval(w.options.Frequency, func() error {
		if err := w.loop(); err != nil {
			log.Println("worker loop iteration had error:", err)
			// never return an error, which would stop the worker
			// instead, we remain hopeful the next iteration will work
		}
		return nil
	})
}

func (w *worker) loop() error {
	log.Println("worker loop starting")
	// we have a list of modules already downloaded to fs
	// we have a list of modules from registry that we want
	// do a diff, finding:
	// - modules we have but do not need anymore
	// - modules we need but to not have yet
	// then prune modules we do not want
	// then DL and save modules we do want
	return nil
}
