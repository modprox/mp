package pokes

import (
	"fmt"

	"github.com/modprox/libmodprox/netservice"
)

type instance = netservice.Instance

type Heartbeat struct {
	Self        instance `json:"self"`
	NumPackages int      `json:"num_packages"`
	NumModules  int      `json:"num_modules"`
}

func (h Heartbeat) String() string {
	return fmt.Sprintf("[%s:%d %d %d]",
		h.Self.Address,
		h.Self.Port,
		h.NumPackages,
		h.NumModules,
	)
}

type StartConfig struct {
	Self       instance `json:"self"`
	Transforms string   `json:"transforms"`
}
