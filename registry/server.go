package registry

import (
	"github.com/modprox/modprox-registry/internal/service"
	"github.com/modprox/modprox-registry/registry/config"
)

func Start(config config.Configuration) {
	service.NewRegistry(config).Run()
}
