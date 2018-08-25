package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/modprox/libmodprox/webutil"
	"github.com/modprox/modprox-registry/internal/repositories"
	"github.com/modprox/modprox-registry/internal/repositories/repository"
)

const (
	get  = http.MethodGet
	post = http.MethodPost
)

func NewRouter(store repositories.Store) http.Handler {
	router := mux.NewRouter()

	router.Handle("/v1/registry/sources/list", registryList(store)).Methods(get)
	router.Handle("/v1/registry/sources/new", registryAdd(store)).Methods(post)

	return router
}

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

		itemsAdded, err := store.Append(wantToAdd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		msg := fmt.Sprintf("added %d submitted items", itemsAdded)
		webutil.WriteJSON(w, msg)
	}
}
