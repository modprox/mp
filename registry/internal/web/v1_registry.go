package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/modprox/libmodprox/coordinates"
	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/libmodprox/webutil"
	"github.com/modprox/modprox-registry/registry/internal/data"
)

func registryAdd(store data.Store) http.HandlerFunc {
	log := loggy.New("registry-add-api")

	return func(w http.ResponseWriter, r *http.Request) {
		log.Tracef("adding to the registry")

		var wantToAdd []coordinates.Module

		if err := json.NewDecoder(r.Body).Decode(&wantToAdd); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		modulesAdded, err := store.InsertModules(wantToAdd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		msg := fmt.Sprintf("added %d new modules", modulesAdded)
		webutil.WriteJSON(w, msg)
	}
}
