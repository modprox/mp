package service

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/modprox/libmodprox/clients/registry"
	"github.com/modprox/libmodprox/clients/zips"
	"github.com/modprox/libmodprox/netservice"
	"github.com/modprox/libmodprox/upstream"
	"github.com/modprox/modprox-proxy/internal/modules/background"
	"github.com/modprox/modprox-proxy/internal/modules/store"
	"github.com/modprox/modprox-proxy/internal/status/heartbeat"
	"github.com/modprox/modprox-proxy/internal/web"
)

type initer func(*Proxy) error

func initIndex(p *Proxy) error {
	// this is the same as the store path for now,
	// because for MVP the index is just another view
	// of the filesystem where the modules are
	//
	// later, we could keep the index in memory if performance
	// is lacking reading the filesystem all the time
	storePath := p.config.ModuleStorage.Path
	if storePath == "" {
		return errors.New("module_storage.path is required")
	}

	p.index = store.NewIndex(store.IndexOptions{
		Directory: storePath,
	})

	return nil
}

func initStore(p *Proxy) error {
	storePath := p.config.ModuleStorage.Path
	if storePath == "" {
		return errors.New("module_storage.path is required")
	}

	tmpPath := p.config.ModuleStorage.Tmp

	p.store = store.NewStore(store.Options{
		Directory:    storePath,
		TmpDirectory: tmpPath,
	})

	return nil
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
	})

	return nil
}

func initZipsClient(p *Proxy) error {
	httpClient := zips.NewHTTPClient(
		zips.HTTPOptions{
			Timeout: 1 * time.Minute,
		},
	)
	p.zipsClient = zips.NewClient(httpClient)
	return nil
}

func initRegistryReloader(p *Proxy) error {
	reloadFreqS := time.Duration(p.config.Registry.PollFrequencyS) * time.Second
	p.reloader = background.NewReloader(
		background.Options{
			Frequency: reloadFreqS,
		},
		p.registryClient,
		p.index,
		p.store,
		upstream.NewResolver(
			initTransforms(p)...,
		),
		p.zipsClient,
	)
	p.reloader.Start()
	return nil
}

func initTransforms(p *Proxy) []upstream.Transform {
	transforms := make([]upstream.Transform, 0, 1)
	transforms = append(transforms, initGoGetTransform(p))
	transforms = append(transforms, initStaticRedirectTransforms(p)...)
	transforms = append(transforms, initSetPathTransform(p))
	transforms = append(transforms, initHeaderTransforms(p)...)
	return transforms
}

func initGoGetTransform(p *Proxy) upstream.Transform {
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

func initHeartbeatSender(p *Proxy) error {
	hostname, err := os.Hostname()
	if err != nil {
		return errors.Wrap(err, "failed to acquire hostname")
	}

	sender := heartbeat.NewSender(heartbeat.Options{
		Timeout:    30 * time.Second,
		Registries: p.config.Registry.Instances,
		Self: netservice.Instance{
			Address: hostname,
			Port:    p.config.APIServer.Port,
		},
	})

	looper := heartbeat.NewLooper(
		10*time.Second,
		sender,
	)

	go looper.Loop()

	return nil
}

func initWebServer(p *Proxy) error {
	router := web.NewRouter(p.index, p.store)

	listenOn := fmt.Sprintf(
		"%s:%d",
		p.config.APIServer.BindAddress,
		p.config.APIServer.Port,
	)

	go func(h http.Handler) {
		if err := http.ListenAndServe(listenOn, h); err != nil {
			p.log.Errorf("failed to listen and serve forever")
			p.log.Errorf("caused by: %v", err)
			os.Exit(1)
		}
	}(router)

	return nil
}
