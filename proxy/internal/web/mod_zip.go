package web

import (
	"net/http"

	"github.com/modprox/mp/pkg/metrics/stats"

	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/proxy/internal/modules/store"
	"github.com/modprox/mp/proxy/internal/web/output"
)

type moduleZip struct {
	store   store.ZipStore
	emitter stats.Sender
	log     loggy.Logger
}

func modZip(store store.ZipStore, emitter stats.Sender) http.Handler {
	return &moduleZip{
		store:   store,
		emitter: emitter,
		log:     loggy.New("mod-zip"),
	}
}

// e.g. GET http://localhost:9000/github.com/shoenig/toolkit/@v/v1.0.0.zip

func (h *moduleZip) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mod, err := modInfoFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.emitter.Count("mod-zip-bad-request", 1)
		return
	}

	h.log.Infof("serving request for .zip file of %s", mod)

	zipBlob, err := h.store.GetZip(mod)
	if err != nil {
		h.log.Warnf("failed to get zip file of %s, %v", mod, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		h.emitter.Count("mod-zip-not-found", 1)
		return
	}

	h.log.Infof("sending zip which is %d bytes", len(zipBlob))
	output.WriteZip(w, zipBlob)
	h.emitter.Count("mod-zip-ok", 1)
}
