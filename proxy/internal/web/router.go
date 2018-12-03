package web

import (
	"net/http"
	"strings"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/gorilla/mux"

	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/pkg/webutil"
	"github.com/modprox/mp/proxy/internal/modules/store"
)

const (
	get  = http.MethodGet
	post = http.MethodPost
)

func NewRouter(
	middles []webutil.Middleware,
	index store.Index,
	store store.ZipStore,
	statter statsd.Statter,
) http.Handler {

	router := mux.NewRouter()

	// e.g. GET  http://localhost:9000/github.com/shoenig/toolkit/@v/v1.0.0.info
	// e.g. GET  http://localhost:9000/github.com/shoenig/toolkit/@v.list
	// e.g. POST http://localhost:9000/github.com/shoenig/toolkit/@v/v1.0.0.rm

	router.PathPrefix("/").Handler(modList(index, statter)).MatcherFunc(suffix("list")).Methods(get)
	router.PathPrefix("/").Handler(modInfo(index, statter)).MatcherFunc(suffix(".info")).Methods(get)
	router.PathPrefix("/").Handler(modFile(index, statter)).MatcherFunc(suffix(".mod")).Methods(get)
	router.PathPrefix("/").Handler(modZip(store, statter)).MatcherFunc(suffix(".zip")).Methods(get)
	router.PathPrefix("/").Handler(modRM(index, store, statter)).MatcherFunc(suffix(".rm")).Methods(post)
	router.PathPrefix("/").HandlerFunc(notFound(statter))
	return webutil.Chain(router, middles...)
}

func suffix(s string) mux.MatcherFunc {
	log := loggy.New("suffix-match")

	return func(r *http.Request, rm *mux.RouteMatch) bool {
		match := strings.HasSuffix(r.URL.Path, s)
		log.Tracef("request from %s matches suffix %q: %t", r.RemoteAddr, s, match)
		return match
	}
}

func notFound(statter statsd.Statter) http.HandlerFunc {
	log := loggy.New("not-found")
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof("request from %s wanted %q which is not found", r.RemoteAddr, r.URL.String())
		http.Error(w, "not found", http.StatusNotFound)
		statter.Inc("path-not-found", 1, 1)
	}
}
