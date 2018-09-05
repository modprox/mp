package proxy

import (
	"github.com/modprox/modprox-proxy/internal/service"
	"github.com/modprox/modprox-proxy/proxy/config"
)

// Start a proxy service instance parameterized by the given Configuration.
func Start(configuration config.Configuration) {
	service.NewProxy(configuration).Run()
}
