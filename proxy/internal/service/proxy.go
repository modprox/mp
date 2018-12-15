package service

import (
	"os"

	"github.com/modprox/mp/proxy/internal/modules/get"

	"github.com/modprox/mp/pkg/clients/registry"
	"github.com/modprox/mp/pkg/clients/zips"
	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/pkg/metrics/stats"
	"github.com/modprox/mp/pkg/webutil"
	"github.com/modprox/mp/proxy/config"
	"github.com/modprox/mp/proxy/internal/modules/bg"
	"github.com/modprox/mp/proxy/internal/modules/store"
	"github.com/modprox/mp/proxy/internal/problems"
)

type Proxy struct {
	config         config.Configuration
	middles        []webutil.Middleware
	emitter        stats.Sender
	index          store.Index
	store          store.ZipStore
	registryClient registry.Client
	zipsClient     zips.Client
	downloader     get.Downloader
	bgWorker       bg.Worker
	dlTracker      problems.Tracker
	log            loggy.Logger
}

func NewProxy(configuration config.Configuration) *Proxy {
	p := &Proxy{
		config: configuration,
		log:    loggy.New("proxy-service"),
	}

	for _, f := range []initer{
		initSender,
		initTrackers,
		initIndex,
		initStore,
		initRegistryClient,
		initZipsClient,
		initBGWorker,
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
