package service

import (
	"net/http"
	"time"

	"github.com/modprox/libmodprox/clients/registry"
	"github.com/modprox/libmodprox/clients/zips"
	"github.com/modprox/libmodprox/upstream"
	"github.com/modprox/modprox-proxy/internal/modules/background"
	"github.com/modprox/modprox-proxy/internal/modules/store"
	"github.com/modprox/modprox-proxy/internal/web"
)

type initer func(*Proxy) error

func initIndex(p *Proxy) error {
	p.index = store.NewIndex(store.IndexOptions{
		Directory: "/tmp/foo",
	})
	return nil
}

func initStore(p *Proxy) error {
	p.store = store.NewStore(store.Options{
		Directory: "/tmp/foo",
	})
	return nil
}

func initRegistryClient(p *Proxy) error {
	p.registryClient = registry.NewClient(registry.Options{
		Timeout:   10 * time.Second,
		Addresses: p.config.Registries,
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

func initReloader(p *Proxy) error {
	pollFreq := time.Duration(p.config.PollRegFreq) * time.Second
	p.reloader = background.NewReloader(
		background.Options{
			Frequency: pollFreq,
		},
		p.registryClient,
		p.store,
		upstream.NewResolver(
			upstream.NewRedirectTransform("example", "code.example.com"),
			upstream.NewSetPathTransform(nil),
		),
		p.zipsClient,
	)
	p.reloader.Start()
	return nil
}

func initWebserver(p *Proxy) error {
	go func(r http.Handler) {
		if err := http.ListenAndServe(":9000", r); err != nil {
			p.log.Errorf("failed to listen and serve forever %v", err)
			panic(err)
		}
	}(web.NewRouter(p.index))
	return nil
}
