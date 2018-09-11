package proxy

import (
	"github.com/modprox/mp/proxy/config"
	"github.com/modprox/mp/proxy/internal/service"
)

// Start a proxy service instance parameterized by the given Configuration.
func Start(configuration config.Configuration) {
	service.NewProxy(configuration).Run()
}
