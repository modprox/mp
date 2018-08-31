package web

import (
	"net/http"

	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/modprox-proxy/internal/modules/store"
	"github.com/modprox/modprox-proxy/internal/web/output"
)

type moduleInfo struct {
	log   loggy.Logger
	index store.Index
}

func modInfo(index store.Index) http.Handler {
	return &moduleInfo{
		index: index,
		log:   loggy.New("mod-info"),
	}
}

// e.g. GET http://localhost:9000/github.com/shoenig/toolkit/@v/v1.0.0.info

func (h *moduleInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mod, err := modInfoFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.log.Infof("serving request for .info of: %s", mod)

	revInfo, err := h.index.Info(mod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content := revInfo.String()
	output.Write(w, output.JSON, content)
}
