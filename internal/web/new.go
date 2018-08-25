package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/modprox/modprox-registry/internal/repositories"
	"github.com/modprox/modprox-registry/static"
)

type newPage struct {
}

type newHandler struct {
	html  *template.Template
	store repositories.Store
}

func newNewHandler(store repositories.Store) http.Handler {
	html := static.MustParseTemplates(
		"static/html/layout.html",
		"static/html/navbar.html",
		"static/html/new.html",
	)

	return &newHandler{
		html:  html,
		store: store,
	}
}

func (h *newHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("[web] new page", r.Method)

	if err := h.html.Execute(w, nil); err != nil {
		log.Panic("failed to serve new page:", err)
	}
}
