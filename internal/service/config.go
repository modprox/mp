package service

import (
	"encoding/json"

	"github.com/modprox/libmodprox/netutil"
)

type Configuration struct {
	Registries  []netutil.Service `json:"registries"`
	PollRegFreq int               `json:"registry_poll_frequency_s"`
}

func (c Configuration) String() string {
	bs, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(bs)
}
