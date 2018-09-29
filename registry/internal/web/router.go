package web

import (
	"net/http"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/gorilla/mux"
	"github.com/shoenig/petrify/v4"

	"github.com/modprox/mp/pkg/webutil"
	"github.com/modprox/mp/registry/internal/data"
	"github.com/modprox/mp/registry/static"
)

const (
	get  = http.MethodGet
	post = http.MethodPost
)

func NewRouter(
	middleAPI []webutil.Middleware,
	middleUI []webutil.Middleware,
	store data.Store,
	statter statsd.Statter,
) http.Handler {

	// 1) a router onto which sub-routers will be mounted
	router := http.NewServeMux()

	// 2) a static files handler for statics
	router.Handle("/static/", routeStatics(http.FileServer(&petrify.AssetFS{
		Asset:     static.Asset,
		AssetDir:  static.AssetDir,
		AssetInfo: static.AssetInfo,
		Prefix:    "static",
	})))

	// 3) an API handler, not CSRF protected
	router.Handle("/v1/", routeAPI(middleAPI, store, statter))

	// 4) a webUI handler, is CSRF protected
	router.Handle("/", routeWebUI(middleUI, store, statter))

	return router
}

func routeStatics(files http.Handler) http.Handler {
	sub := mux.NewRouter()
	sub.Handle("/static/css/{*}", http.StripPrefix("/static/", files)).Methods(get)
	sub.Handle("/static/img/{*}", http.StripPrefix("/static/", files)).Methods(get)
	return sub
}

func routeAPI(middles []webutil.Middleware, store data.Store, stats statsd.Statter) http.Handler {
	sub := mux.NewRouter()
	sub.Handle("/v1/registry/sources/list", newRegistryList(store, stats)).Methods(get, post)
	sub.Handle("/v1/registry/sources/new", registryAdd(store, stats)).Methods(post)
	sub.Handle("/v1/proxy/heartbeat", newHeartbeatHandler(store, stats)).Methods(post)
	sub.Handle("/v1/proxy/configuration", newStartupHandler(store, stats)).Methods(post)
	return webutil.Chain(sub, middles...)
}

func routeWebUI(middles []webutil.Middleware, store data.Store, stats statsd.Statter) http.Handler {
	sub := mux.NewRouter()
	sub.Handle("/mods/new", newAddHandler(store, stats)).Methods(get, post)
	sub.Handle("/mods/list", newModsListHandler(store, stats)).Methods(get)
	sub.Handle("/mods/show", newShowHandler(store, stats)).Methods(get)
	sub.Handle("/configure/about", newAboutHandler(stats)).Methods(get)
	sub.Handle("/configure/blocks", newBlocksHandler(stats)).Methods(get)
	sub.Handle("/", newHomeHandler(store, stats)).Methods(get, post)
	return webutil.Chain(sub, middles...)
}
