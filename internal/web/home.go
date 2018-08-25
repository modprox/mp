package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/modprox/modprox-registry/internal/repositories"
	"github.com/modprox/modprox-registry/internal/repositories/repository"
	"github.com/modprox/modprox-registry/static"
)

type homepage struct {
	Modules []repository.Module
}

func newHomepageHandler(store repositories.Store) http.Handler {
	html := static.MustParseTemplates(
		"static/html/layout.html",
		"static/html/navbar.html",
		"static/html/home.html",
	)
	return &homepageHandler{
		html:  html,
		store: store,
	}
}

type homepageHandler struct {
	html  *template.Template
	store repositories.Store
}

func (h *homepageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("[web] serving up the homepage")

	if err := h.html.Execute(w, nil); err != nil {
		log.Panicf("failed to execute homepage template: %v", err)
	}
}
