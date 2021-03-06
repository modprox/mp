package service

import (
	"os"

	"gophers.dev/pkgs/loggy"

	"oss.indeed.com/go/modprox/pkg/clients/zips"
	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/registry/config"
	"oss.indeed.com/go/modprox/registry/internal/data"
)

type Registry struct {
	config      config.Configuration
	store       data.Store
	emitter     stats.Sender
	log         loggy.Logger
	history     string
	proxyClient zips.ProxyClient
}

func NewRegistry(config config.Configuration) *Registry {
	r := &Registry{
		config: config,
		log:    loggy.New("registry-service"),
	}

	for _, f := range []initer{
		initSender,
		initStore,
		initProxyPrune,
		initHistory,
		initProxyClient,
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
