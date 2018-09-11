package registry

import (
	"github.com/modprox/modprox-registry/registry/config"
	"github.com/modprox/modprox-registry/registry/internal/service"
)

func Start(config config.Configuration) {
	service.NewRegistry(config).Run()
}
