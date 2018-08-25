package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/modprox/libmodprox/webutil"
	"github.com/modprox/modprox-registry/internal/repositories"
	"github.com/modprox/modprox-registry/internal/repositories/repository"
)

func registryList(store repositories.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("[api] registry list endpoint")

		repos, err := store.ListSources()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		webutil.WriteJSON(w, repos)
	}
}

func registryAdd(store repositories.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var wantToAdd []repository.Module

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
