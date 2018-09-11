package proxy

import (
	"github.com/modprox/modprox-proxy/proxy/config"
	"github.com/modprox/modprox-proxy/proxy/internal/service"
)

// Start a proxy service instance parameterized by the given Configuration.
func Start(configuration config.Configuration) {
	service.NewProxy(configuration).Run()
}
