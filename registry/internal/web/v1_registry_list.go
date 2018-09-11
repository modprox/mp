package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/modprox/libmodprox/clients/registry"
	"github.com/modprox/libmodprox/coordinates"
	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/libmodprox/webutil"
	"github.com/modprox/modprox-registry/registry/internal/data"
)

type registryList struct {
	store data.Store
	log   loggy.Logger
}

func newRegistryList(store data.Store) http.Handler {
	return &registryList{
		store: store,
		log:   loggy.New("registry-list-api"),
	}
}

func (h *registryList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var send toSend

	switch r.Method {
	case http.MethodGet:
		send = h.get(w, r)
	case http.MethodPost:
		send = h.post(w, r)
	}

	if send.err != nil {
		h.log.Errorf("failed to serve request: %v", send.err)
		http.Error(w, send.err.Error(), send.code)
		return
	}

	response := registry.ReqModsResp{
		Mods: send.mods,
	}

	webutil.WriteJSON(w, response)
}

type toSend struct {
	err  error
	code int
	mods []coordinates.SerialModule
}

func (h *registryList) get(w http.ResponseWriter, r *http.Request) toSend {
	h.log.Tracef("listing entire contents of registry")
	modules, err := h.store.ListModules()
	if err != nil {
		return toSend{
			err:  err,
			code: http.StatusInternalServerError,
			mods: nil,
		}
	}
	return toSend{
		err:  nil,
		code: http.StatusOK,
		mods: modules,
	}
}

func (h *registryList) post(w http.ResponseWriter, r *http.Request) toSend {
	h.log.Tracef("listing optimized contents of registry")

	// read the body of the incoming request
	var inbound registry.ReqMods
	if err := json.NewDecoder(r.Body).Decode(&inbound); err != nil {
		return toSend{
			err:  err,
			code: http.StatusBadRequest,
			mods: nil,
		}
	}

	fmt.Println("ranges:", inbound.IDs)

	ids, err := h.store.ListModuleIDs()
	if err != nil {
		return toSend{
			err:  err,
			code: http.StatusInternalServerError,
			mods: nil,
		}
	}

	fmt.Println("ids:", ids)

	// compare that with the modules in the registry

	// return a list of the difference
	neededIDs := inListButNotRange(ids, inbound.IDs)

	fmt.Println("needed ids:", neededIDs)

	needed, err := h.store.ListModulesByIDs(neededIDs)
	if err != nil {
		return toSend{
			err:  err,
			code: http.StatusInternalServerError,
			mods: nil,
		}
	}

	fmt.Println("needed mods:", needed)

	return toSend{
		err:  nil,
		code: http.StatusOK,
		mods: needed,
	}
}

// this could be optimized doing a kind of skipping merge, but for
// now the O(n) scan should be okay
func inListButNotRange(ids []int64, ranges coordinates.RangeIDs) []int64 {
	var neededIDs []int64
	for _, id := range ids {
		needsID := true
		for _, r := range ranges {
			if inRange(id, r) {
				needsID = false
				break // move on to next id
			}
		}
		if needsID {
			neededIDs = append(neededIDs, id)
		}
	}
	return neededIDs
}

func inRange(i int64, rangeID coordinates.RangeID) bool {
	left := rangeID[0]
	right := rangeID[1]
	return i >= left && i <= right
}
