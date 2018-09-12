package web

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/shoenig/petrify/v4"

	"github.com/modprox/mp/registry/config"
	"github.com/modprox/mp/registry/internal/data"
	"github.com/modprox/mp/registry/static"
)

const (
	get  = http.MethodGet
	post = http.MethodPost
)

func NewRouter(
	store data.Store,
	csrfConfig config.CSRF,
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
	router.Handle("/v1/", routeAPI(store))

	// 4) a webUI handler, is CSRF protected
	router.Handle("/", routeWebUI(csrfConfig, store))

	return router
}

func routeStatics(files http.Handler) http.Handler {
	sub := mux.NewRouter()
	sub.Handle("/static/css/{*}", http.StripPrefix("/static/", files)).Methods(get)
	sub.Handle("/static/img/{*}", http.StripPrefix("/static/", files)).Methods(get)
	return sub
}

func routeAPI(store data.Store) http.Handler {
	sub := mux.NewRouter()
	sub.Handle("/v1/registry/sources/list", newRegistryList(store)).Methods(get, post)
	sub.Handle("/v1/registry/sources/new", registryAdd(store)).Methods(post)
	sub.Handle("/v1/proxy/heartbeat", newHeartbeatHandler(store)).Methods(post)
	sub.Handle("/v1/proxy/configuration", newStartupHandler(store)).Methods(post)
	return sub
}

func routeWebUI(csrfConfig config.CSRF, store data.Store) http.Handler {
	sub := mux.NewRouter()
	sub.Handle("/new", newAddHandler(store)).Methods(get, post)
	sub.Handle("/configure/redirects", newRedirectsHandler(store)).Methods(get)
	sub.Handle("/", newHomeHandler(store)).Methods(get, post)
	return chain(sub,
		[]middleware{csrf.Protect(
			// the key is used to generate csrf tokens to hand
			// out on html form loads
			[]byte(csrfConfig.AuthenticationKey),

			// CSRF cookies are https-only normally, so for development
			//// mode make sure the csrf package knows we are using http
			csrf.Secure(!csrfConfig.DevelopmentMode),
		)}...,
	)
}

type middleware func(http.Handler) http.Handler

//  chain recursively chains middleware together
func chain(h http.Handler, m ...middleware) http.Handler {
	if len(m) == 0 {
		return h
	}
	return m[0](chain(h, m[1:cap(m)]...))
}
