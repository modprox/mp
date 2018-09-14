package service

import (
	"os"

	"github.com/cactus/go-statsd-client/statsd"

	"github.com/modprox/mp/pkg/clients/registry"
	"github.com/modprox/mp/pkg/clients/zips"
	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/proxy/config"
	"github.com/modprox/mp/proxy/internal/modules/background"
	"github.com/modprox/mp/proxy/internal/modules/store"
)

type Proxy struct {
	config         config.Configuration
	statter        statsd.Statter
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
		initStatter,
		initIndex,
		initStore,
		initRegistryClient,
		initZipsClient,
		initRegistryReloader,
		initHeartbeatSender,
		initStartupConfigSender,
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
