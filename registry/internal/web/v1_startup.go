package web

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"oss.indeed.com/go/modprox/pkg/metrics/stats"

	"oss.indeed.com/go/modprox/pkg/clients/payloads"
	"oss.indeed.com/go/modprox/pkg/loggy"
	"oss.indeed.com/go/modprox/proxy/config"
	"oss.indeed.com/go/modprox/registry/internal/data"
)

type startupHandler struct {
	store   data.Store
	emitter stats.Sender
	log     loggy.Logger
}

func newStartupHandler(store data.Store, emitter stats.Sender) http.Handler {
	return &startupHandler{
		store:   store,
		emitter: emitter,
		log:     loggy.New("startup-config-handler"),
	}
}

func (h *startupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Tracef("receiving startup configuration from proxy")

	code, msg, err := h.post(r)
	if err != nil {
		h.log.Errorf("failed to accept startup configuration from %s, %v", r.RemoteAddr, err)
		http.Error(w, msg, code)
		h.emitter.Count("api-proxy-start-config-error", 1)
		return
	}

	h.log.Tracef("accepted startup configuration from %s", r.RemoteAddr)
	_, _ = io.WriteString(w, "ok")
	h.emitter.Count("api-proxy-start-config-ok", 1)
}

func (h *startupHandler) post(r *http.Request) (int, string, error) {
	// proxy should probably send an Instance to identify itself
	var configuration payloads.Configuration

	if err := json.NewDecoder(r.Body).Decode(&configuration); err != nil {
		return http.StatusBadRequest, "failed to decode request", err
	}

	if err := checkConfiguration(configuration); err != nil {
		return http.StatusBadRequest, "configuration is not valid", err
	}

	if err := h.store.SetStartConfig(configuration); err != nil {
		return http.StatusInternalServerError, "failed to save configuration", err
	}

	return http.StatusOK, "ok", nil
}

func checkConfiguration(configuration payloads.Configuration) error {
	switch {
	case configuration.Storage == config.Storage{}:
		return errors.New("storage configuration cannot be empty")
	case len(configuration.Registry.Instances) == 0:
		return errors.New("registries configuration cannot be empty")
	}
	return nil
}
