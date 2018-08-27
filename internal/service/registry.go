package service

import (
	"os"

	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/modprox-registry/internal/data"
)

type Registry struct {
	config Configuration
	store  data.Store
	log    loggy.Logger
}

func NewRegistry(config Configuration) *Registry {
	r := &Registry{
		config: config,
		log:    loggy.New("registry-service"),
	}

	for _, f := range []initer{
		initStore,
		initWebserver,
	} {
		if err := f(r); err != nil {
			r.log.Errorf("failed to initialize registry: %v", err)
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
