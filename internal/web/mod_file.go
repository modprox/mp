package web

import (
	"net/http"

	"github.com/modprox/libmodprox/loggy"
)

type moduleFile struct {
	log loggy.Logger
}

func newModuleFile() http.Handler {
	return &moduleFile{
		log: loggy.New("mod-file"),
	}
}

func (h *moduleFile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Infof("serving request for module file")
	return
}
