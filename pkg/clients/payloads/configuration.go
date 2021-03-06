package payloads

import (
	"encoding/json"

	"github.com/jinzhu/copier"

	"oss.indeed.com/go/modprox/pkg/netservice"
	"oss.indeed.com/go/modprox/pkg/setup"
	"oss.indeed.com/go/modprox/proxy/config"
)

// Configuration of a proxy instance when it starts up that is sent
// to the registry.
type Configuration struct {
	Self            netservice.Instance   `json:"self"`
	DiskStorage     config.Storage        `json:"disk_storage,omitempty"`
	DatabaseStorage setup.PersistentStore `json:"database_storage,omitempty"`
	Registry        config.Registry       `json:"registry"`
	Transforms      config.Transforms     `json:"transforms"`
}

func (c Configuration) Texts() (string, string, string, error) {

	storageText, err := json.Marshal(c.DiskStorage)
	if err != nil {
		return "", "", "", err
	}

	registriesText, err := json.Marshal(c.Registry)
	if err != nil {
		return "", "", "", err
	}

	// hide the values of the headers, which may contain secrets
	var t2 config.Transforms

	if err := copier.Copy(&t2, &c.Transforms); err != nil {
		return "", "", "", err
	}

	for _, transform := range t2.DomainHeaders {
		for key := range transform.Headers {
			transform.Headers[key] = "********"
		}
	}

	transformsText, err := json.Marshal(t2)
	if err != nil {
		return "", "", "", err
	}

	return string(storageText), string(registriesText), string(transformsText), nil
}
