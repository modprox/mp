package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/modprox/modprox-registry/internal/repositories"
	"github.com/modprox/modprox-registry/internal/repositories/repository"
	"github.com/modprox/libmodprox/webutil"
)

const (
	get  = http.MethodGet
	post = http.MethodPost
)

func NewRouter(store repositories.Store) http.Handler {
	router := mux.NewRouter()

	router.Handle("/v1/registry/list", registryList(store)).Methods(get)
	router.Handle("/v1/registry/append", registryAdd(store)).Methods(post)

	return router
}

func registryList(store repositories.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repos, err := store.List()
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

		if err := store.Append(wantToAdd); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		itemsAdded := len(wantToAdd)
		msg := fmt.Sprintf("added %d submitted items", itemsAdded)
		webutil.WriteJSON(w, msg)
	}
}
