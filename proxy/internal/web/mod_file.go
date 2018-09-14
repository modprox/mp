package web

import (
	"net/http"

	"github.com/cactus/go-statsd-client/statsd"

	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/proxy/internal/modules/store"
	"github.com/modprox/mp/proxy/internal/web/output"
)

type moduleFile struct {
	index   store.Index
	statter statsd.Statter
	log     loggy.Logger
}

func modFile(index store.Index, statter statsd.Statter) http.Handler {
	return &moduleFile{
		index:   index,
		statter: statter,
		log:     loggy.New("mod-file"),
	}
}

func (h *moduleFile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mod, err := modInfoFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.statter.Inc("mod-file-bad-request", 1, 1)
		return
	}
	h.log.Infof("serving request for go.mod file of %s", mod)

	modFile, err := h.index.Mod(mod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		h.statter.Inc("mod-file-not-found", 1, 1)
		return
	}

	output.Write(w, output.Text, modFile)
	h.statter.Inc("mod-file-ok", 1, 1)
}
