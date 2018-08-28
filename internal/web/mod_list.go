package web

import (
	"net/http"

	"github.com/modprox/libmodprox/loggy"
)

type moduleList struct {
	log loggy.Logger
}

func newModuleList() http.Handler {
	return &moduleList{
		log: loggy.New("mod-list"),
	}
}

func (h *moduleList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Infof("serving request for listing")
	return
}
