package service

import (
	"github.com/modprox/libmodprox/loggy"

	"github.com/modprox/modprox-proxy/internal/modules/background"
	"github.com/modprox/modprox-proxy/internal/modules/store"
)

type Proxy struct {
	config   Configuration
	store    store.Store
	reloader background.Reloader
	log      loggy.Logger
}

func NewProxy(config Configuration) *Proxy {
	p := &Proxy{
		config: config,
		log:    loggy.New("proxy-service"),
	}

	for _, i := range []initer{
		initStore,
		initReloader,
		initWebserver,
	} {
		if err := i(p); err != nil {
			p.log.Errorf("failed to initialize proxy: %v", err)
			panic(err)
		}
	}

	return p
}

func (p *Proxy) Run() {
	select {
	// intentionally left blank
	}
}
