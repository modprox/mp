package web

import (
	"html/template"
	"net/http"

	"github.com/cactus/go-statsd-client/statsd"

	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/registry/static"
)

type aboutHandler struct {
	html    *template.Template
	statter statsd.Statter
	log     loggy.Logger
}

func newAboutHandler(statter statsd.Statter) http.Handler {
	html := static.MustParseTemplates(
		"static/html/layout.html",
		"static/html/navbar.html",
		"static/html/about.html",
	)
	return &aboutHandler{
		html:    html,
		statter: statter,
		log:     loggy.New("about-handler"),
	}
}

func (h *aboutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.html.Execute(w, nil); err != nil {
		h.log.Errorf("failed to execute about template: %v", err)
		return
	}

	h.statter.Inc("ui-about-ok", 1, 1)
}
