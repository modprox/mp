package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/libmodprox/netservice"
	"github.com/modprox/libmodprox/pokes"
	"github.com/modprox/libmodprox/webutil"
	"github.com/modprox/modprox-registry/internal/data"
)

type heartbeatHandler struct {
	store data.Store
	log   loggy.Logger
}

func newHeartbeatHandler(store data.Store) http.Handler {
	return &heartbeatHandler{
		store: store,
		log:   loggy.New("heartbeat-update-handler"),
	}
}

func (h *heartbeatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Tracef("receiving heartbeat update")

	var (
		code     int
		response string
		from     netservice.Instance
		err      error
	)

	code, response, from, err = h.post(r)

	if err != nil {
		h.log.Errorf("failed to accept heartbeat from %s, %v", from, err)
		http.Error(w, response, code)
		return
	}

	h.log.Tracef("accepted heartbeat from %s", from)
	webutil.WriteJSON(w, response)
}

func (h *heartbeatHandler) post(r *http.Request) (int, string, netservice.Instance, error) {
	var from netservice.Instance
	var heartbeat pokes.Heartbeat
	if err := json.NewDecoder(r.Body).Decode(&heartbeat); err != nil {
		return http.StatusBadRequest, "failed to decode request", from, err
	}

	if err := checkHeartbeat(heartbeat); err != nil {
		return http.StatusBadRequest, "heartbeat is not valid", from, err
	}

	from = heartbeat.Self

	if err := h.store.SetHeartbeat(heartbeat); err != nil {
		return http.StatusInternalServerError, "failed to save heartbeat", from, err
	}

	return http.StatusOK, "ok", from, nil
}

func checkHeartbeat(heartbeat pokes.Heartbeat) error {
	switch {
	case heartbeat.Self.Address == "":
		return errors.New("heartbeat address cannot be empty")
	case heartbeat.Self.Port <= 0:
		return errors.New("heartbeat port must be positive")
	case heartbeat.NumPackages < 0:
		return errors.New("heartbeat num_packages must be non-negative")
	case heartbeat.NumModules < 0:
		return errors.New("heartbeat num_modules must be non-negative")
	}
	return nil
}
