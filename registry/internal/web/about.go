package web

import (
	"html/template"
	"net/http"

	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/pkg/metrics/stats"
	"github.com/modprox/mp/registry/static"
)

type aboutHandler struct {
	html    *template.Template
	emitter stats.Sender
	log     loggy.Logger
}

func newAboutHandler(emitter stats.Sender) http.Handler {
	html := static.MustParseTemplates(
		"static/html/layout.html",
		"static/html/navbar.html",
		"static/html/about.html",
	)
	return &aboutHandler{
		html:    html,
		emitter: emitter,
		log:     loggy.New("about-handler"),
	}
}

func (h *aboutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.html.Execute(w, nil); err != nil {
		h.log.Errorf("failed to execute about template: %v", err)
		return
	}

	h.emitter.Count("ui-about-ok", 1)
}
