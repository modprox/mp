package web

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"go.gophers.dev/pkgs/loggy"

	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/pkg/webutil"
	"oss.indeed.com/go/modprox/proxy/internal/modules/store"
	"oss.indeed.com/go/modprox/proxy/internal/problems"
)

const (
	get  = http.MethodGet
	post = http.MethodPost
)

func NewRouter(
	middles []webutil.Middleware,
	index store.Index,
	store store.ZipStore,
	emitter stats.Sender,
	dlProblems problems.Tracker,
	history string,
) http.Handler {

	router := mux.NewRouter()

	// mod operations
	//
	// e.g. GET  http://localhost:9000/github.com/example/toolkit/@v/v1.0.0.info
	// e.g. GET  http://localhost:9000/github.com/example/toolkit/@v.list
	// e.g. POST http://localhost:9000/github.com/example/toolkit/@v/v1.0.0.rm
	router.PathPrefix("/").Handler(modList(index, emitter)).MatcherFunc(suffix("list")).Methods(get)
	router.PathPrefix("/").Handler(modInfo(index, emitter)).MatcherFunc(suffix(".info")).Methods(get)
	router.PathPrefix("/").Handler(modFile(index, emitter)).MatcherFunc(suffix(".mod")).Methods(get)
	router.PathPrefix("/").Handler(modZip(store, emitter)).MatcherFunc(suffix(".zip")).Methods(get)
	router.PathPrefix("/").Handler(modRM(index, store, emitter)).MatcherFunc(suffix(".rm")).Methods(post)

	// metadata about this app
	router.PathPrefix("/history").Handler(appHistory(emitter, history)).Methods(get)

	// api operations
	//
	router.PathPrefix("/v1/problems/downloads").Handler(newDownloadProblems(dlProblems, emitter)).Methods(get)

	// default behavior (404)
	router.PathPrefix("/").HandlerFunc(notFound(emitter))

	// force middleware
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

func notFound(emitter stats.Sender) http.HandlerFunc {
	log := loggy.New("not-found")
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof("request from %s wanted %q which is not found", r.RemoteAddr, r.URL.String())
		http.Error(w, "not found", http.StatusNotFound)
		emitter.Count("path-not-found", 1)
	}
}
