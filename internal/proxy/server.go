package proxy

import (
	"github.com/modprox/modprox-proxy/internal/proxy/config"
	"github.com/modprox/modprox-proxy/internal/service"
)

func Start(configuration config.Configuration) {
	service.NewProxy(configuration).Run()
}
