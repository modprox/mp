package web

import (
	"net/http"

	"github.com/gorilla/mux"

	petrify "go.gophers.dev/cmds/petrify/v5"

	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/pkg/webutil"
	"oss.indeed.com/go/modprox/registry/internal/data"
	"oss.indeed.com/go/modprox/registry/static"
)

const (
	get  = http.MethodGet
	post = http.MethodPost
)

func NewRouter(
	middleAPI []webutil.Middleware,
	middleUI []webutil.Middleware,
	store data.Store,
	emitter stats.Sender,
	history string,
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
	router.Handle("/v1/", routeAPI(middleAPI, store, emitter))

	// 4) a webUI handler, is CSRF protected
	router.Handle("/", routeWebUI(middleUI, store, emitter, history))

	return router
}

func routeStatics(files http.Handler) http.Handler {
	sub := mux.NewRouter()
	sub.Handle("/static/css/{*}", http.StripPrefix("/static/", files)).Methods(get)
	sub.Handle("/static/img/{*}", http.StripPrefix("/static/", files)).Methods(get)
	return sub
}

func routeAPI(middles []webutil.Middleware, store data.Store, emitter stats.Sender) http.Handler {
	sub := mux.NewRouter()
	sub.Handle("/v1/registry/sources/list", newRegistryList(store, emitter)).Methods(get, post)
	sub.Handle("/v1/registry/sources/new", registryAdd(store, emitter)).Methods(post)
	sub.Handle("/v1/proxy/heartbeat", newHeartbeatHandler(store, emitter)).Methods(post)
	sub.Handle("/v1/proxy/configuration", newStartupHandler(store, emitter)).Methods(post)
	return webutil.Chain(sub, middles...)
}

func routeWebUI(middles []webutil.Middleware, store data.Store, emitter stats.Sender, history string) http.Handler {
	sub := mux.NewRouter()
	sub.Handle("/mods/new", newAddHandler(store, emitter)).Methods(get, post)
	sub.Handle("/mods/list", newModsListHandler(store, emitter)).Methods(get)
	sub.Handle("/mods/show", newShowHandler(store, emitter)).Methods(get, post)
	sub.Handle("/mods/find", newFindHandler(emitter)).Methods(get, post)
	sub.Handle("/configure/about", newAboutHandler(emitter)).Methods(get)
	sub.Handle("/configure/blocks", newBlocksHandler(emitter)).Methods(get)
	sub.Handle("/history", newHistoryHandler(emitter, history)).Methods(get)
	sub.Handle("/", newHomeHandler(store, emitter)).Methods(get, post)
	return webutil.Chain(sub, middles...)
}
