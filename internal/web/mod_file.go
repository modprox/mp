package web

import (
	"net/http"

	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/modprox-proxy/internal/modules/store"
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
	h.log.Infof("serving request for module file")
	return
}
