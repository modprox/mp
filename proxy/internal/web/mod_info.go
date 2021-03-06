package web

import (
	"net/http"

	"gophers.dev/pkgs/loggy"

	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/proxy/internal/modules/store"
	"oss.indeed.com/go/modprox/proxy/internal/web/output"
)

type moduleInfo struct {
	log     loggy.Logger
	emitter stats.Sender
	index   store.Index
}

func modInfo(index store.Index, emitter stats.Sender) http.Handler {
	return &moduleInfo{
		index:   index,
		emitter: emitter,
		log:     loggy.New("mod-info"),
	}
}

// e.g. GET http://localhost:9000/github.com/example/toolkit/@v/v1.0.0.info

func (h *moduleInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mod, err := modInfoFromPath(r.URL.Path)
	if err != nil {
		h.log.Warnf("bad request for info: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.emitter.Count("mod-info-bad-request", 1)
		return
	}

	h.log.Infof("serving request for .info of: %s", mod)

	revInfo, err := h.index.Info(mod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		h.emitter.Count("mod-info-not-found", 1)
		return
	}

	content := revInfo.String()
	output.Write(w, output.JSON, content)
	h.emitter.Count("mod-info-ok", 1)
}
