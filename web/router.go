package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", nothing)

	return router
}

func nothing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello\n"))
}
