package proxy

import (
	"oss.indeed.com/go/modprox/proxy/config"
	"oss.indeed.com/go/modprox/proxy/internal/service"
)

// Start a proxy service instance parameterized by the given Configuration.
func Start(configuration config.Configuration) {
	service.NewProxy(configuration).Run()
}
