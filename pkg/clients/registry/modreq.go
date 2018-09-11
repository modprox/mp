package registry

import "github.com/modprox/mp/pkg/coordinates"

// ReqMods is the data sent from the proxy to the registry when
// requesting from the registry a list of modules that the proxy
// is in need of downloading to its local data-store. When making
// the request, the proxy sends a list of ranges of serial IDs of
// the modules it already has contained in its data-store. That way
// the registry can compare that list of ranges with the set of all
// modules that are registered, and reply with a list of modules that
// only contains modules the proxy needs to download.
type ReqMods struct {
	IDs coordinates.RangeIDs `json:"ids"`
}

// ReqModsResp is the response sent from the registry to the proxy
// when the proxy requests a list of which modules it needs to download.
// There is no guarantee the proxy will not have already downloaded
// some of the modules, but given the implementation it should be pretty
// well optimized to not include duplicates of what is in the proxy store.
type ReqModsResp struct {
	Mods []coordinates.SerialModule `json:"serials"`
}
