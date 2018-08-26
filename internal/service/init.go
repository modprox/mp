package service

import (
	"log"
	"net/http"
	"time"

	"github.com/modprox/modprox-proxy/internal/modules/background"
	"github.com/modprox/modprox-proxy/internal/modules/store"
	"github.com/modprox/modprox-proxy/internal/web"
)

type initer func(*Proxy) error

func initStore(p *Proxy) error {
	p.store = store.NewStore(store.Options{
		Directory: "/tmp/foo",
	})
	return nil
}

func initReloader(p *Proxy) error {
	pollFreq := time.Duration(p.config.PollRegFreq) * time.Second
	p.reloader = background.NewReloader(
		background.Options{
			Frequency: pollFreq,
		},
		p.store,
	)
	p.reloader.Start()
	return nil
}

func initWebserver(p *Proxy) error {
	go func(r http.Handler) {
		if err := http.ListenAndServe(":9000", r); err != nil {
			log.Fatalf("failed to listen and serve forever %v", err)
		}
	}(web.NewRouter())
	return nil
}
