package web

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/shoenig/petrify/v4"

	"github.com/modprox/modprox-registry/internal/data"
	"github.com/modprox/modprox-registry/registry/config"
	"github.com/modprox/modprox-registry/static"
)

const (
	get  = http.MethodGet
	post = http.MethodPost
)

func NewRouter(
	store data.Store,
	csrfConfig config.CSRF,
) http.Handler {

	// 1) a router onto which subrouters will be mounted
	router := mux.NewRouter()

	// 2) a static files handler for statics
	routeStatics(router, http.FileServer(&petrify.AssetFS{
		Asset:     static.Asset,
		AssetDir:  static.AssetDir,
		AssetInfo: static.AssetInfo,
		Prefix:    "static",
	}))

	// 3) an API handler, not CSRF protected
	routeAPI(router, store)

	// 4) a webUI handler, is CSRF protected
	routeWebUI(router, csrfConfig, store)

	return router
}

func routeStatics(router *mux.Router, files http.Handler) http.Handler {
	sub := router.PathPrefix("/").Subrouter()
	sub.Handle("/static/css/{*}", http.StripPrefix("/static/", files)).Methods(get)
	sub.Handle("/static/img/{*}", http.StripPrefix("/static/", files)).Methods(get)
	return sub
}

func routeAPI(router *mux.Router, store data.Store) http.Handler {
	sub := router.PathPrefix("/v1").Subrouter()
	sub.Handle("/registry/sources/list", registryList(store)).Methods(get)
	sub.Handle("/registry/sources/new", registryAdd(store)).Methods(post)
	sub.Handle("/heartbeat/update", newHeartbeatHandler(store)).Methods(post)
	return sub
}

func routeWebUI(router *mux.Router, csrfConfig config.CSRF, store data.Store) http.Handler {
	sub := router.PathPrefix("/").Subrouter()
	sub.Handle("/new", newAddHandler(store)).Methods(get, post)
	sub.Handle("/configure/redirects", newRedirectsHandler(store)).Methods(get)
	sub.Handle("/", newHomeHandler(store)).Methods(get)

	middlewares := []middleware{
		csrf.Protect(
			// the key is used to generate csrf tokens to hand
			// out on html form loads
			[]byte(csrfConfig.AuthenticationKey),

			// CSRF cookies are https-only normally, so for development
			//// mode make sure the csrf package knows we are using http
			csrf.Secure(!csrfConfig.DevelopmentMode),
		),
	}

	return chain(sub, middlewares...)
}

type middleware func(http.Handler) http.Handler

//  chain recursively chains middleware together
func chain(h http.Handler, m ...middleware) http.Handler {
	if len(m) == 0 {
		return h
	}
	return m[0](chain(h, m[1:cap(m)]...))
}
