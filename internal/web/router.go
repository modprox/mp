package web

import (
	"net/http"

	"github.com/modprox/libmodprox/loggy"

	"github.com/gorilla/mux"
)

func NewRouter() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", nothing)

	return router
}

func nothing(w http.ResponseWriter, r *http.Request) {
	log := loggy.New("nothing")

	log.Tracef("doing nothing in this handler!")

	w.Write([]byte("hello\n"))
}
