package web

import (
	"net/http"
	"strings"

	"go.gophers.dev/pkgs/loggy"

	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/proxy/internal/modules/store"
	"oss.indeed.com/go/modprox/proxy/internal/web/output"
)

type moduleList struct {
	index   store.Index
	emitter stats.Sender
	log     loggy.Logger
}

func modList(index store.Index, emitter stats.Sender) http.Handler {
	return &moduleList{
		index:   index,
		emitter: emitter,
		log:     loggy.New("mod-list"),
	}
}

func (h *moduleList) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	module, err := moduleFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.log.Infof("serving request for list module: %s", module)

	listing, err := h.index.Versions(module)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		h.emitter.Count("mod-list-not-found", 1)
		return
	}

	output.Write(w, output.Text, formatList(listing))
	h.emitter.Count("mod-list-ok", 1)
}

func formatList(list []string) string {
	var sb strings.Builder
	for _, version := range list {
		sb.WriteString(version)
		sb.WriteString("\n")
	}
	return sb.String()
}
