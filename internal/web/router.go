package web

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/modprox/modprox-registry/internal/repositories"
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
