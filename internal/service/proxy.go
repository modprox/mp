package service

import (
	"log"

	"github.com/modprox/modprox-proxy/internal/modules/background"
	"github.com/modprox/modprox-proxy/internal/modules/store"
)

type Proxy struct {
	config   Configuration
	store    store.Store
	reloader background.Reloader
}

func NewProxy(config Configuration) *Proxy {
	p := &Proxy{config: config}

	for _, i := range []initer{
		initStore,
		initReloader,
		initWebserver,
	} {
		if err := i(p); err != nil {
			log.Fatal("failed to initialize proxy:", err)
		}
	}

	return p
}

func (p *Proxy) Run() {
	select {
	// intentionally left blank
	}
}
