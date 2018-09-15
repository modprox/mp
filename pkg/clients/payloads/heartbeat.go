package payloads

import (
	"fmt"

	"github.com/modprox/mp/pkg/netservice"
)

type instance = netservice.Instance

// Heartbeat of a proxy that is sent periodically to the registry.
type Heartbeat struct {
	Self        instance `json:"self"`
	NumModules  int      `json:"num_modules"`
	NumVersions int      `json:"num_versions"`
}

func (h Heartbeat) String() string {
	return fmt.Sprintf("[%s:%d %d %d]",
		h.Self.Address,
		h.Self.Port,
		h.NumModules,
		h.NumVersions,
	)
}
