package registry

import (
	"github.com/modprox/mp/registry/config"
	"github.com/modprox/mp/registry/internal/service"
)

func Start(config config.Configuration) {
	service.NewRegistry(config).Run()
}
