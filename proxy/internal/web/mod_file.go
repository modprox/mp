package web

import (
	"net/http"

	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/proxy/internal/modules/store"
	"github.com/modprox/mp/proxy/internal/web/output"
)

type moduleFile struct {
	index store.Index
	log   loggy.Logger
}

func modFile(index store.Index) http.Handler {
	return &moduleFile{
		index: index,
		log:   loggy.New("mod-file"),
	}
}

func (h *moduleFile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mod, err := modInfoFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.log.Infof("serving request for go.mod file of %s", mod)

	modFile, err := h.index.Mod(mod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	output.Write(w, output.Text, modFile)
}
