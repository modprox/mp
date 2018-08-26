package service

import (
	"github.com/modprox/libmodprox/configutil"
	"github.com/modprox/libmodprox/netutil"
)

type Configuration struct {
	Registries  []netutil.Service `json:"registries"`
	PollRegFreq int               `json:"registry_poll_frequency_s"`
}

func (c Configuration) String() string {
	return configutil.Format(c)
}
