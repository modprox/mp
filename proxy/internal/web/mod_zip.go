package web

import (
	"net/http"

	"github.com/cactus/go-statsd-client/statsd"

	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/proxy/internal/modules/store"
	"github.com/modprox/mp/proxy/internal/web/output"
)

type moduleZip struct {
	store   store.ZipStore
	statter statsd.Statter
	log     loggy.Logger
}

func modZip(store store.ZipStore, statter statsd.Statter) http.Handler {
	return &moduleZip{
		store:   store,
		statter: statter,
		log:     loggy.New("mod-zip"),
	}
}

// e.g. GET http://localhost:9000/github.com/shoenig/toolkit/@v/v1.0.0.zip

func (h *moduleZip) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mod, err := modInfoFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.statter.Inc("mod-zip-bad-request", 1, 1)
		return
	}

	h.log.Infof("serving request for .zip file of %s", mod)

	zipBlob, err := h.store.GetZip(mod)
	if err != nil {
		h.log.Warnf("failed to get zip file of %s, %v", mod, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		h.statter.Inc("mod-zip-not-found", 1, 1)
		return
	}

	h.log.Infof("sending zip which is %d bytes", len(zipBlob))
	output.WriteZip(w, zipBlob)
	h.statter.Inc("mod-zip-ok", 1, 1)
}
