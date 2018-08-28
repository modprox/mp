package web

import (
	"net/http"

	"github.com/modprox/libmodprox/loggy"
)

type moduleZip struct {
	log loggy.Logger
}

func newModuleZip() http.Handler {
	return &moduleZip{
		log: loggy.New("mod-zip"),
	}
}

func (h *moduleZip) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Infof("serving request for zip file")
}
