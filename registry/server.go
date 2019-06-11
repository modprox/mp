package registry

import (
	"oss.indeed.com/go/modprox/registry/config"
	"oss.indeed.com/go/modprox/registry/internal/service"
)

func Start(config config.Configuration) {
	service.NewRegistry(config).Run()
}
