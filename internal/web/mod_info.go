package web

import (
	"net/http"

	"github.com/modprox/libmodprox/loggy"
)

type moduleInfo struct {
	log loggy.Logger
}

func newModuleInfo() http.Handler {
	return &moduleInfo{
		log: loggy.New("mod-info"),
	}
}

func (h *moduleInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Infof("serving request for info")
}
