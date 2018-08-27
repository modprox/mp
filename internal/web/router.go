package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shoenig/petrify/v4"

	"github.com/modprox/modprox-registry/internal/data"
	"github.com/modprox/modprox-registry/static"
)

const (
	get  = http.MethodGet
	post = http.MethodPost
)

func NewRouter(store data.Store) http.Handler {
	router := mux.NewRouter()

	staticFiles := http.FileServer(&petrify.AssetFS{
		Asset:     static.Asset,
		AssetDir:  static.AssetDir,
		AssetInfo: static.AssetInfo,
		Prefix:    "static",
	})

	// v1 API
	router.Handle("/v1/registry/sources/list", registryList(store)).Methods(get)
	router.Handle("/v1/registry/sources/new", registryAdd(store)).Methods(post)

	// website
	router.Handle("/static/css/{*}", http.StripPrefix("/static/", staticFiles)).Methods(get)
	// router.Handle("/static/imgs/{*}", http.StripPrefix("/static/", staticFiles)).Methods(get)

	router.Handle("/new", newAddHandler(store)).Methods(get, post)
	router.Handle("/configure/redirects", newRedirectsHandler(store)).Methods(get)
	router.Handle("/", newHomeHandler(store)).Methods(get)
	return router
}
