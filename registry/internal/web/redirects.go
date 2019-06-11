package web

import (
	"html/template"
	"net/http"

	"oss.indeed.com/go/modprox/pkg/loggy"
	"oss.indeed.com/go/modprox/registry/internal/data"
	"oss.indeed.com/go/modprox/registry/static"
)

type redirectsHandler struct {
	html  *template.Template
	store data.Store
	log   loggy.Logger
}

func newRedirectsHandler(store data.Store) http.Handler {
	html := static.MustParseTemplates(
		"static/html/layout.html",
		"static/html/navbar.html",
		"static/html/redirects.html",
	)

	return &redirectsHandler{
		html:  html,
		store: store,
		log:   loggy.New("redirects-handler"),
	}
}

func (h *redirectsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Tracef("loaded page %v", r.Method)

	if err := h.html.Execute(w, nil); err != nil {
		h.log.Errorf("failed to execute redirects template: %v", err)
		return
	}
}
