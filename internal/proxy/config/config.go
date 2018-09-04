package config

import (
	"github.com/modprox/libmodprox/configutil"
	"github.com/modprox/libmodprox/netservice"
)

type Configuration struct {
	APIServer     APIServer `json:"api_server"`
	Registry      Registry  `json:"registry"`
	ModuleStorage Storage   `json:"module_storage"`
}

func (c Configuration) String() string {
	return configutil.Format(c)
}

type APIServer struct {
	TLS struct {
		Enabled     bool   `json:"enabled"`
		Certificate string `json:"certificate"`
		Key         string `json:"key"`
	} `json:"tls"`
	BindAddress string `json:"bind_address"`
	Port        int    `json:"port"`
}

type Storage struct {
	Path string `json:"path"`
}

type instances = []netservice.Instance

type Registry struct {
	PollFrequencyS  int       `json:"poll_frequency_s"`
	RequestTimeoutS int       `json:"request_timeout_s"`
	Instances       instances `json:"instances"`
}
