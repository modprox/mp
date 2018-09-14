package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cactus/go-statsd-client/statsd"

	"github.com/modprox/mp/pkg/coordinates"
	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/pkg/webutil"
	"github.com/modprox/mp/registry/internal/data"
)

func registryAdd(store data.Store, statter statsd.Statter) http.HandlerFunc {
	log := loggy.New("registry-add-api")

	return func(w http.ResponseWriter, r *http.Request) {
		log.Tracef("adding to the registry")

		var wantToAdd []coordinates.Module

		if err := json.NewDecoder(r.Body).Decode(&wantToAdd); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			statter.Inc("api-addmod-bad-request", 1, 1)
			return
		}

		modulesAdded, err := store.InsertModules(wantToAdd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			statter.Inc("api-addmod-error", 1, 1)
			return
		}

		msg := fmt.Sprintf("added %d new modules", modulesAdded)
		webutil.WriteJSON(w, msg)
		statter.Inc("api-addmod-ok", 1, 1)
	}
}
