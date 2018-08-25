package service

import (
	"encoding/json"
)

type Configuration struct {
	Registries  []Registry `json:"registries"`
	PollRegFreq int        `json:"registry_poll_frequency_s"`
}

func (c Configuration) String() string {
	bs, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(bs)
}

type Registry struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}
