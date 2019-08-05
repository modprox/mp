package service

import (
	"os"

	"gophers.dev/pkgs/loggy"

	"oss.indeed.com/go/modprox/pkg/clients/registry"
	"oss.indeed.com/go/modprox/pkg/clients/zips"
	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/pkg/webutil"
	"oss.indeed.com/go/modprox/proxy/config"
	"oss.indeed.com/go/modprox/proxy/internal/modules/bg"
	"oss.indeed.com/go/modprox/proxy/internal/modules/get"
	"oss.indeed.com/go/modprox/proxy/internal/modules/store"
	"oss.indeed.com/go/modprox/proxy/internal/problems"
)

type Proxy struct {
	config         config.Configuration
	middles        []webutil.Middleware
	emitter        stats.Sender
	index          store.Index
	store          store.ZipStore
	registryClient registry.Client
	proxyClient    zips.ProxyClient
	upstreamClient zips.UpstreamClient
	downloader     get.Downloader
	bgWorker       bg.Worker
	dlTracker      problems.Tracker
	log            loggy.Logger
	history        string
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
		initZipClients,
		initBGWorker,
		initHeartbeatSender,
		initStartupConfigSender,
		initHistory,
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
