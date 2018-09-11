package service

import (
	"os"

	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/registry/config"
	"github.com/modprox/mp/registry/internal/data"
)

type Registry struct {
	config config.Configuration
	store  data.Store
	log    loggy.Logger
}

func NewRegistry(config config.Configuration) *Registry {
	r := &Registry{
		config: config,
		log:    loggy.New("registry-service"),
	}

	for _, f := range []initer{
		initStore,
		initWebServer,
	} {
		if err := f(r); err != nil {
			r.log.Errorf("cannot startup: failed to initialize registry")
			r.log.Errorf("caused by: %v", err)
			os.Exit(1)
		}
	}

	return r
}

func (r *Registry) Run() {
	select {
	//intentionally left blank
	}
}
