package get

import (
	"bytes"
	"encoding/json"

	"oss.indeed.com/go/modprox/pkg/clients/registry"
	"oss.indeed.com/go/modprox/pkg/coordinates"
	"oss.indeed.com/go/modprox/pkg/loggy"
	"oss.indeed.com/go/modprox/proxy/internal/modules/store"
)

// Range is an alias of coordinates.RangeIDs for brevity.
type Ranges = coordinates.RangeIDs

// RegistryAPI is used to issue API request from the registry
type RegistryAPI interface {
	ModulesNeeded(Ranges) ([]coordinates.SerialModule, error)
}

type registryAPI struct {
	registryClient registry.Client
	index          store.Index
	log            loggy.Logger
}

func NewRegistryAPI(
	registryClient registry.Client,
	index store.Index,
) RegistryAPI {
	return &registryAPI{
		registryClient: registryClient,
		index:          index,
		log:            loggy.New("registryAPI"),
	}
}

func (r *registryAPI) ModulesNeeded(excludeIDs Ranges) ([]coordinates.SerialModule, error) {
	ids, err := r.index.IDs()
	if err != nil {
		return nil, err
	}

	rm := registry.ReqMods{
		IDs: ids,
	}

	bs, err := json.Marshal(rm)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(bs)

	var buf bytes.Buffer
	if err := r.registryClient.Post("/v1/registry/sources/list", reader, &buf); err != nil {
		return nil, err
	}

	r2 := bytes.NewReader(buf.Bytes())

	var response registry.ReqModsResp
	if err := json.NewDecoder(r2).Decode(&response); err != nil {
		return nil, err
	}

	return response.Mods, nil
}
