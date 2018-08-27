package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/libmodprox/repository"
	"github.com/modprox/libmodprox/webutil"
	"github.com/modprox/modprox-registry/internal/repositories"
)

func registryList(store repositories.Store) http.HandlerFunc {
	log := loggy.New("registry-list-api")

	return func(w http.ResponseWriter, r *http.Request) {
		log.Tracef("listing contents of registry")

		repos, err := store.ListSources()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		webutil.WriteJSON(w, repos)
	}
}

func registryAdd(store repositories.Store) http.HandlerFunc {
	log := loggy.New("registry-add-api")

	return func(w http.ResponseWriter, r *http.Request) {
		log.Tracef("adding to the registry")

		var wantToAdd []repository.ModInfo

		if err := json.NewDecoder(r.Body).Decode(&wantToAdd); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sourcesAdded, tagsAdded, err := store.Add(wantToAdd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		msg := fmt.Sprintf("added %d tags across %d sources", tagsAdded, sourcesAdded)
		webutil.WriteJSON(w, msg)
	}
}
