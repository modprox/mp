package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"gophers.dev/pkgs/loggy"

	"oss.indeed.com/go/modprox/pkg/clients/payloads"
	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/pkg/netservice"
	"oss.indeed.com/go/modprox/pkg/webutil"
	"oss.indeed.com/go/modprox/registry/internal/data"
)

type heartbeatHandler struct {
	store   data.Store
	emitter stats.Sender
	log     loggy.Logger
}

func newHeartbeatHandler(store data.Store, emitter stats.Sender) http.Handler {
	return &heartbeatHandler{
		store:   store,
		emitter: emitter,
		log:     loggy.New("heartbeat-update-handler"),
	}
}

func (h *heartbeatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code, response, from, err := h.post(r)

	if err != nil {
		h.log.Errorf("failed to accept heartbeat from %s, %v", from, err)
		http.Error(w, response, code)
		h.emitter.Count("api-heartbeat-unaccepted", 1)
		return
	}

	h.log.Tracef("accepted heartbeat from %s", from)
	webutil.WriteJSON(w, response)
	h.emitter.Count("api-heartbeat-accepted", 1)
}

func (h *heartbeatHandler) post(r *http.Request) (int, string, netservice.Instance, error) {
	var from netservice.Instance
	var heartbeat payloads.Heartbeat
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

func checkHeartbeat(heartbeat payloads.Heartbeat) error {
	switch {
	case heartbeat.Self.Address == "":
		return errors.New("heartbeat address cannot be empty")
	case heartbeat.Self.Port <= 0:
		return errors.New("heartbeat port must be positive")
	case heartbeat.NumModules < 0:
		return errors.New("heartbeat num_packages must be non-negative")
	case heartbeat.NumVersions < 0:
		return errors.New("heartbeat num_modules must be non-negative")
	}
	return nil
}
