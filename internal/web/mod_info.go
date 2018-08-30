package web

import (
	"net/http"

	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/modprox-proxy/internal/modules/store"
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
	h.log.Infof("serving request for info")
}
