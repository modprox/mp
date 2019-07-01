package web

import (
	"net/http"

	"go.gophers.dev/pkgs/loggy"

	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/proxy/internal/web/output"
)

type applicationHistory struct {
	log     loggy.Logger
	history string
	emitter stats.Sender
}

func appHistory(emitter stats.Sender, history string) http.Handler {
	return &applicationHistory{
		emitter: emitter,
		history: history,
		log:     loggy.New("app-history"),
	}
}

// e.g. GET http://localhost:9000/history

func (h *applicationHistory) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	output.Write(w, output.Text, h.history)
	h.emitter.Count("app-history-ok", 1)
}
