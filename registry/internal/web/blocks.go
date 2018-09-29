package web

import (
	"html/template"
	"net/http"

	"github.com/cactus/go-statsd-client/statsd"

	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/registry/static"
)

type blocksHandler struct {
	html    *template.Template
	statter statsd.Statter
	log     loggy.Logger
}

func newBlocksHandler(statter statsd.Statter) http.Handler {
	html := static.MustParseTemplates(
		"static/html/layout.html",
		"static/html/navbar.html",
		"static/html/blocks.html",
	)

	return &blocksHandler{
		html:    html,
		statter: statter,
		log:     loggy.New("blocks-handler"),
	}
}

func (h *blocksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.html.Execute(w, nil); err != nil {
		h.log.Errorf("failed to execute blocks template: %v", err)
		return
	}

	h.statter.Inc("ui-blocks-ok", 1, 1)
}
