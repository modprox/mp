package payloads

import "github.com/modprox/mp/proxy/config"

// Configuration of a proxy instance when it starts up that is sent
// to the registry.
type Configuration struct {
	Storage    config.Storage    `json:"storage"`
	Registry   config.Registry   `json:"registry"`
	Transforms config.Transforms `json:"transforms"`
}
