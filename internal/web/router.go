package web

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/modprox-proxy/internal/modules/store"
)

const (
	get = http.MethodGet
)

func NewRouter(
	index store.Index,
	store store.ZipStore,
) http.Handler {

	router := mux.NewRouter()

	// e.g. GET http://localhost:9000/github.com/shoenig/toolkit/@v/v1.0.0.info

	router.PathPrefix("/").Handler(modList(index)).MatcherFunc(suffix("list")).Methods(get)
	router.PathPrefix("/").Handler(modInfo(index)).MatcherFunc(suffix(".info")).Methods(get)
	router.PathPrefix("/").Handler(modFile(index)).MatcherFunc(suffix(".mod")).Methods(get)
	router.PathPrefix("/").Handler(modZip(store)).MatcherFunc(suffix(".zip")).Methods(get)
	router.PathPrefix("/").HandlerFunc(notFound())

	return router
}

func suffix(s string) mux.MatcherFunc {
	log := loggy.New("suffix-match")

	return func(r *http.Request, rm *mux.RouteMatch) bool {
		match := strings.HasSuffix(r.URL.Path, s)
		log.Tracef("request from %s matches suffix %q: %t", r.RemoteAddr, s, match)
		return match
	}
}

func notFound() http.HandlerFunc {
	log := loggy.New("not-found")
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof("request from %s wanted %q which is not found", r.RemoteAddr, r.URL.String())
		http.Error(w, "not found", http.StatusNotFound)
	}
}
