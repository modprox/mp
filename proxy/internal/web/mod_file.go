package web

import (
	"net/http"

	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/pkg/metrics/stats"
	"github.com/modprox/mp/proxy/internal/modules/store"
	"github.com/modprox/mp/proxy/internal/web/output"
)

type moduleFile struct {
	index   store.Index
	emitter stats.Sender
	log     loggy.Logger
}

func modFile(index store.Index, emitter stats.Sender) http.Handler {
	return &moduleFile{
		index:   index,
		emitter: emitter,
		log:     loggy.New("mod-file"),
	}
}

func (h *moduleFile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mod, err := modInfoFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.emitter.Count("mod-file-bad-request", 1)
		return
	}
	h.log.Infof("serving request for go.mod file of %s", mod)

	modFile, err := h.index.Mod(mod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		h.emitter.Count("mod-file-not-found", 1)
		return
	}

	output.Write(w, output.Text, modFile)
	h.emitter.Count("mod-file-ok", 1)
}
