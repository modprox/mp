package web

import (
	"net/http"

	"go.gophers.dev/pkgs/loggy"

	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/proxy/internal/problems"
	"oss.indeed.com/go/modprox/proxy/internal/web/output"
)

type downloadProblems struct {
	dlTracker problems.Tracker
	emitter   stats.Sender
	log       loggy.Logger
}

func newDownloadProblems(dlTracker problems.Tracker, emitter stats.Sender) http.Handler {
	return &downloadProblems{
		dlTracker: dlTracker,
		emitter:   emitter,
		log:       loggy.New("download-problems"),
	}
}

func (h *downloadProblems) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.emitter.Count("api-download-problems", 1)

	dlProblems := h.dlTracker.Problems()
	h.log.Tracef("reporting %d download problems", len(dlProblems))

	output.WriteJSON(w, dlProblems)
}
