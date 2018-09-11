package service

import (
	"os"

	"github.com/modprox/libmodprox/clients/registry"
	"github.com/modprox/libmodprox/clients/zips"
	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/modprox-proxy/proxy/config"
	"github.com/modprox/modprox-proxy/proxy/internal/modules/background"
	"github.com/modprox/modprox-proxy/proxy/internal/modules/store"
)

type Proxy struct {
	config         config.Configuration
	index          store.Index
	store          store.ZipStore
	registryClient registry.Client
	zipsClient     zips.Client
	reloader       background.Reloader
	log            loggy.Logger
}

func NewProxy(configuration config.Configuration) *Proxy {
	p := &Proxy{
		config: configuration,
		log:    loggy.New("proxy-service"),
	}

	for _, f := range []initer{
		initIndex,
		initStore,
		initRegistryClient,
		initZipsClient,
		initRegistryReloader,
		initHeartbeatSender,
		initWebServer,
	} {
		if err := f(p); err != nil {
			p.log.Errorf("failed to initialize proxy: %v", err)
			os.Exit(1)
		}
	}

	return p
}

func (p *Proxy) Run() {
	select {
	// intentionally left blank
	}
}
