package web

import (
	"net/http"

	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/modprox-registry/internal/data"
)

type heartbeatHandler struct {
	store data.Store
	log   loggy.Logger
}

func newHeartbeatHandler(store data.Store) http.Handler {
	return &heartbeatHandler{
		store: store,
		log:   loggy.New("heartbeat-handler"),
	}
}

func (h *heartbeatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Tracef("receiving hearbeat update")
}
