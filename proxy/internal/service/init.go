package service

import (
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"

	"oss.indeed.com/go/modprox/pkg/clients/payloads"
	"oss.indeed.com/go/modprox/pkg/clients/registry"
	"oss.indeed.com/go/modprox/pkg/clients/zips"
	"oss.indeed.com/go/modprox/pkg/history"
	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/pkg/netservice"
	"oss.indeed.com/go/modprox/pkg/setup"
	"oss.indeed.com/go/modprox/pkg/upstream"
	"oss.indeed.com/go/modprox/pkg/webutil"
	"oss.indeed.com/go/modprox/proxy/internal/modules/bg"
	"oss.indeed.com/go/modprox/proxy/internal/modules/get"
	"oss.indeed.com/go/modprox/proxy/internal/modules/store"
	"oss.indeed.com/go/modprox/proxy/internal/problems"
	"oss.indeed.com/go/modprox/proxy/internal/status/heartbeat"
	"oss.indeed.com/go/modprox/proxy/internal/status/startup"
	"oss.indeed.com/go/modprox/proxy/internal/web"
)

type initer func(*Proxy) error

func initSender(p *Proxy) error {
	cfg := p.config.Statsd.Agent
	if cfg.Port == 0 || cfg.Address == "" {
		p.emitter = stats.Discard()
		p.log.Warnf("stats emitter is set to discard client - no metrics will be reported")
		return nil
	}

	emitter, err := stats.New(stats.Proxy, p.config.Statsd)
	if err != nil {
		return err
	}
	p.emitter = emitter
	return nil
}

func initTrackers(p *Proxy) error {
	dlTracker := problems.New("downloads")
	p.dlTracker = dlTracker
	return nil
}

func initIndex(p *Proxy) error {
	if p.config.ModuleStorage == nil && p.config.ModuleDBStorage == nil {
		return errors.New("configs must be specified for either file or db module index")
	} else if p.config.ModuleStorage != nil && p.config.ModuleDBStorage != nil {
		return errors.New("configs must be specified for either file or db module index (but not both)")
	}

	if p.config.ModuleStorage != nil {
		indexPath := p.config.ModuleStorage.IndexPath
		index, err := store.NewIndex(store.IndexOptions{
			Directory: indexPath,
		})
		if err != nil {
			return errors.WithStack(err)
		}
		p.index = index
	} else if p.config.ModuleDBStorage != nil {
		_, dsn, err := dbStorageDSN(p, p.config.ModuleDBStorage)
		if err != nil {
			return errors.WithStack(err)
		}
		p.index, err = store.Connect(dsn, p.emitter)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func initStore(p *Proxy) error {
	if p.config.ModuleStorage == nil && p.config.ModuleDBStorage == nil {
		return errors.New("configs must be specified for either file or db module index")
	} else if p.config.ModuleStorage != nil && p.config.ModuleDBStorage != nil {
		return errors.New("configs must be specified for either file or db module index (but not both)")
	}

	if p.config.ModuleStorage != nil {
		storePath := p.config.ModuleStorage.DataPath
		if storePath == "" {
			return errors.New("module_storage.path is required")
		}

		tmpPath := p.config.ModuleStorage.TmpPath
		p.store = store.NewStore(store.Options{
			Directory:    storePath,
			TmpDirectory: tmpPath,
		}, p.emitter)
	} else if p.config.ModuleDBStorage != nil {
		_, dsn, err := dbStorageDSN(p, p.config.ModuleDBStorage)
		if err != nil {
			return errors.WithStack(err)
		}
		p.store, err = store.Connect(dsn, p.emitter)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func dbStorageDSN(p *Proxy, conf *setup.PersistentStore) (string, setup.DSN, error) {
	kind, dsn, err := conf.DSN()
	if err != nil {
		return "", setup.DSN{}, errors.WithStack(err)
	}
	p.log.Infof("using database of kind: %q", kind)
	p.log.Infof("database dsn: %s", dsn)

	return kind, dsn, nil
}

func initRegistryClient(p *Proxy) error {
	clientTimeout := p.config.Registry.RequestTimeoutS
	if clientTimeout <= 0 {
		return errors.Errorf(
			"registry.request_timeout_s must be > 0, got %d",
			clientTimeout,
		)
	}

	instances := p.config.Registry.Instances
	if len(instances) <= 2 {
		p.log.Warnf(
			"at least 2 registry instances recommended, got %d",
			len(instances),
		)
	}

	p.registryClient = registry.NewClient(registry.Options{
		Timeout:   time.Duration(clientTimeout) * time.Second,
		Instances: p.config.Registry.Instances,
		APIKey:    p.config.Registry.APIKey,
	})

	return nil
}

func initZipClients(p *Proxy) error {
	// create a proxy zip client
	p.proxyClient = zips.NewProxyClient(
		zips.ProxyClientOptions{
			Protocol: p.config.ZipProxy.Protocol,
			BaseURL:  p.config.ZipProxy.BaseURL,
			Timeout:  1 * time.Minute,
		},
	)

	// create an upstream zip client
	httpClient := zips.NewHTTPClient(
		zips.HTTPOptions{
			Timeout: 1 * time.Minute,
		},
	)
	p.upstreamClient = zips.NewUpstreamClient(httpClient)

	return nil
}

func initBGWorker(p *Proxy) error {
	reloadFreqS := time.Duration(p.config.Registry.PollFrequencyS) * time.Second
	registryRequester := get.NewRegistryAPI(
		p.registryClient,
		p.index,
	)

	resolver := upstream.NewResolver(
		initTransforms(p)...,
	)

	downloader := get.New(
		p.proxyClient,
		p.upstreamClient,
		resolver,
		p.store,
		p.index,
		p.emitter,
	)

	p.bgWorker = bg.New(
		p.emitter,
		p.dlTracker,
		p.index,
		p.store,
		registryRequester,
		downloader,
	)

	// start the background worker polling the registry
	p.bgWorker.Start(bg.Options{
		Frequency: reloadFreqS,
	})

	return nil
}

func initTransforms(p *Proxy) []upstream.Transform {
	transforms := make([]upstream.Transform, 0, 1)
	transforms = append(transforms, initGoGetTransform(p))
	transforms = append(transforms, initStaticRedirectTransforms(p)...)
	transforms = append(transforms, initSetPathTransform(p))
	transforms = append(transforms, initHeaderTransforms(p)...)
	transforms = append(transforms, initTransportTransforms(p)...)
	return transforms
}

func initGoGetTransform(p *Proxy) upstream.Transform {
	if p.config.Transforms.AutomaticRedirect {
		return upstream.NewAutomaticGoGetTransform()
	}

	goGetDomains := make([]string, 0, len(p.config.Transforms.DomainGoGet))
	for _, domain := range p.config.Transforms.DomainGoGet {
		goGetDomains = append(goGetDomains, domain.Domain)
	}
	return upstream.NewGoGetTransform(goGetDomains)
}

func initStaticRedirectTransforms(p *Proxy) []upstream.Transform {
	transforms := make([]upstream.Transform, 0, len(p.config.Transforms.DomainRedirects))
	for _, domainRedirect := range p.config.Transforms.DomainRedirects {
		transforms = append(transforms, upstream.NewStaticRedirectTransform(
			domainRedirect.Original,
			domainRedirect.Substitution,
		))
	}
	return transforms
}

func initSetPathTransform(p *Proxy) upstream.Transform {
	transforms := make(map[string]upstream.Transform)
	for _, t := range p.config.Transforms.DomainPath {
		transforms[t.Domain] = upstream.NewDomainPathTransform(t.Path)
	}
	return upstream.NewSetPathTransform(transforms)
}

func initHeaderTransforms(p *Proxy) []upstream.Transform {
	transforms := make([]upstream.Transform, 0, len(p.config.Transforms.DomainHeaders))
	for _, t := range p.config.Transforms.DomainHeaders {
		transforms = append(transforms, upstream.NewDomainHeaderTransform(
			t.Domain, t.Headers,
		))
	}
	return transforms
}

func initTransportTransforms(p *Proxy) []upstream.Transform {
	transforms := make([]upstream.Transform, 0, len(p.config.Transforms.DomainTransport))
	for _, t := range p.config.Transforms.DomainTransport {
		transforms = append(transforms, upstream.NewDomainTransportTransform(
			t.Domain, t.Transport,
		))
	}
	return transforms
}

func initHeartbeatSender(p *Proxy) error {
	sender := heartbeat.NewSender(
		netservice.Instance{
			Address: netservice.Hostname(),
			Port:    p.config.APIServer.Port,
		},
		p.registryClient,
		p.emitter,
	)

	looper := heartbeat.NewLooper(
		10*time.Second,
		p.index,
		p.emitter,
		sender,
	)

	go looper.Loop()

	return nil
}

func initStartupConfigSender(p *Proxy) error {
	sender := startup.NewSender(
		p.registryClient,
		30*time.Second,
		p.emitter,
	)
	go func() {
		cfg := payloads.Configuration{
			Self: netservice.Instance{
				Address: netservice.Hostname(),
				Port:    p.config.APIServer.Port,
			},
			Registry:   p.config.Registry,
			Transforms: p.config.Transforms,
		}
		if p.config.ModuleStorage != nil {
			cfg.DiskStorage = *p.config.ModuleStorage
		}
		if p.config.ModuleDBStorage != nil {
			cfg.DatabaseStorage = *p.config.ModuleDBStorage
		}
		_ = sender.Send(cfg)
	}()
	return nil
}

func initWebServer(p *Proxy) error {
	var middles []webutil.Middleware

	mux := web.NewRouter(
		middles,
		p.index,
		p.store,
		p.emitter,
		p.dlTracker,
		p.history,
	)

	server, err := p.config.APIServer.Server(mux)
	if err != nil {
		return err
	}

	go func(h http.Handler) {
		var err error
		if p.config.APIServer.TLS.Enabled {
			err = server.ListenAndServeTLS(
				p.config.APIServer.TLS.Certificate,
				p.config.APIServer.TLS.Key,
			)
		} else {
			err = server.ListenAndServe()
		}

		p.log.Errorf("server stopped serving: %v", err)
		os.Exit(1)
	}(mux)

	return nil
}

func initHistory(p *Proxy) error {
	historyBytes, err := history.Asset("history.txt")
	if err != nil {
		return err
	}

	p.history = string(historyBytes)
	return nil
}
