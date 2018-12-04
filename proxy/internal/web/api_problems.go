package web

import (
	"net/http"

	"github.com/cactus/go-statsd-client/statsd"

	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/proxy/internal/problems"
	"github.com/modprox/mp/proxy/internal/web/output"
)

type downloadProblems struct {
	dlTracker problems.Tracker
	statter   statsd.Statter
	log       loggy.Logger
}

func newDownloadProblems(dlTracker problems.Tracker, statter statsd.Statter) http.Handler {
	return &downloadProblems{
		dlTracker: dlTracker,
		statter:   statter,
		log:       loggy.New("download-problems"),
	}
}

func (h *downloadProblems) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = h.statter.Inc("api-download-problems", 1, 1)

	dlProblems := h.dlTracker.Problems()
	h.log.Tracef("reporting %d download problems", len(dlProblems))

	output.WriteJSON(w, dlProblems)
}
