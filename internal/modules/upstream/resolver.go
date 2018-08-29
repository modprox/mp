package upstream

import (
	"fmt"

	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/libmodprox/repository"
)

// A Resolver provides the
type Resolver interface {
	Resolve(info repository.ModInfo) string
}

func Passthrough() Resolver {
	return &passthroughHTTP{
		log: loggy.New("resolve"),
	}
}

// e.g. https://github.com/shoenig/petrify/archive/v4.0.1.zip

type passthroughHTTP struct {
	log loggy.Logger
}

func (r *passthroughHTTP) Resolve(mod repository.ModInfo) string {
	address := fmt.Sprintf("https://%s/archive/%s.zip", mod.Source, mod.Version)
	r.log.Tracef("given %s, becomes => %s", mod, address)
	return address
}
