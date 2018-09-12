package payloads

import (
	"fmt"

	"github.com/modprox/mp/pkg/netservice"
)

type instance = netservice.Instance

// Heartbeat of a proxy that is sent periodically to the registry.
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
