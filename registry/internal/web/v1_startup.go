package web

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/cactus/go-statsd-client/statsd"

	"github.com/modprox/mp/pkg/clients/payloads"
	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/proxy/config"
	"github.com/modprox/mp/registry/internal/data"
)

type startupHandler struct {
	store   data.Store
	statter statsd.Statter
	log     loggy.Logger
}

func newStartupHandler(store data.Store, statter statsd.Statter) http.Handler {
	return &startupHandler{
		store:   store,
		statter: statter,
		log:     loggy.New("startup-config-handler"),
	}
}

func (h *startupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Tracef("receiving startup configuration from proxy")

	code, msg, err := h.post(r)
	if err != nil {
		h.log.Errorf("failed to accept startup configuration from %s, %v", r.RemoteAddr, err)
		http.Error(w, msg, code)
		h.statter.Inc("api-proxy-start-config-error", 1, 1)
		return
	}

	h.log.Tracef("accepted startup configuration from %s", r.RemoteAddr)
	io.WriteString(w, "ok")
	h.statter.Inc("api-proxy-start-config-ok", 1, 1)
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
