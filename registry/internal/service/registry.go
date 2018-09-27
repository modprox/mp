package service

import (
	"os"

	"github.com/cactus/go-statsd-client/statsd"

	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/registry/config"
	"github.com/modprox/mp/registry/internal/data"
)

type Registry struct {
	config  config.Configuration
	store   data.Store
	statter statsd.Statter
	log     loggy.Logger
}

func NewRegistry(config config.Configuration) *Registry {
	r := &Registry{
		config: config,
		log:    loggy.New("registry-service"),
	}

	for _, f := range []initer{
		initStatter,
		initStore,
		initProxyPrune,
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
