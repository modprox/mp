package service

import (
	"log"
	"net/http"

	"github.com/modprox/modprox-proxy/internal/web"
)

type Proxy struct {
	config Configuration
}

func NewProxy(config Configuration) *Proxy {
	return &Proxy{
		config: config,
	}
}

func (p *Proxy) Start() {
	router := web.NewRouter()
	if err := http.ListenAndServe(":9000", router); err != nil {
		log.Fatalf("failed to listen and serve forever %v", err)
	}
}
