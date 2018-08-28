package web

import (
	"net/http"

	"github.com/modprox/libmodprox/loggy"

	"github.com/gorilla/mux"
)

func NewRouter() http.Handler {
	router := mux.NewRouter()

	router.Handle("/module/{v}/list", newModuleList())
	router.Handle("/module/{v}/version.info", newModuleInfo())
	router.Handle("/module/{v}/version.mod", newModuleFile())
	router.Handle("/module/{v}/version.zip", newModuleZip())

	router.HandleFunc("/", index)

	return router
}

func index(w http.ResponseWriter, r *http.Request) {
	log := loggy.New("404")
	log.Tracef("remote %s attempted to reach unknown path: %s", r.RemoteAddr, r.URL.Path)
	http.Error(w, "nothing to see here", http.StatusNotFound)
}
