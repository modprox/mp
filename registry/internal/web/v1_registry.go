package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gophers.dev/pkgs/loggy"

	"oss.indeed.com/go/modprox/pkg/coordinates"
	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/pkg/webutil"
	"oss.indeed.com/go/modprox/registry/internal/data"
)

func registryAdd(store data.Store, emitter stats.Sender) http.HandlerFunc {
	log := loggy.New("registry-add-api")

	return func(w http.ResponseWriter, r *http.Request) {
		log.Tracef("adding to the registry")

		var wantToAdd []coordinates.Module

		if err := json.NewDecoder(r.Body).Decode(&wantToAdd); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			emitter.Count("api-addmod-bad-request", 1)
			return
		}

		modulesAdded, err := store.InsertModules(wantToAdd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			emitter.Count("api-addmod-error", 1)
			return
		}

		msg := fmt.Sprintf("added %d new modules", modulesAdded)
		webutil.WriteJSON(w, msg)
		emitter.Count("api-addmod-ok", 1)
	}
}
