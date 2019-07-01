package web

import (
	"html/template"
	"net/http"

	"go.gophers.dev/pkgs/loggy"

	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/registry/static"
)

type blocksHandler struct {
	html    *template.Template
	emitter stats.Sender
	log     loggy.Logger
}

func newBlocksHandler(emitter stats.Sender) http.Handler {
	html := static.MustParseTemplates(
		"static/html/layout.html",
		"static/html/navbar.html",
		"static/html/blocks.html",
	)

	return &blocksHandler{
		html:    html,
		emitter: emitter,
		log:     loggy.New("blocks-handler"),
	}
}

func (h *blocksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.html.Execute(w, nil); err != nil {
		h.log.Errorf("failed to execute blocks template: %v", err)
		return
	}

	h.emitter.Count("ui-blocks-ok", 1)
}
