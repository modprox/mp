package payloads

import (
	"fmt"
	"time"

	"oss.indeed.com/go/modprox/pkg/netservice"
)

type instance = netservice.Instance

// Heartbeat of a proxy that is sent periodically to the registry.
type Heartbeat struct {
	Self        instance `json:"self"`
	NumModules  int      `json:"num_modules"`
	NumVersions int      `json:"num_versions"`
	Timestamp   int      `json:"send_time"` // unix timestamp seconds
}

func (hb Heartbeat) String() string {
	return fmt.Sprintf("[%s:%d %d %d]",
		hb.Self.Address,
		hb.Self.Port,
		hb.NumModules,
		hb.NumVersions,
	)
}

func (hb Heartbeat) TimeSince() string {
	t1 := time.Unix(int64(hb.Timestamp), 0)
	dur := time.Since(t1)
	d := dur.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%dh%dm%ds", h, m, s)
}
