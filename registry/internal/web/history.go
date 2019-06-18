package web

import (
	"net/http"

	"oss.indeed.com/go/modprox/pkg/loggy"
	"oss.indeed.com/go/modprox/pkg/metrics/stats"
)

type historyHandler struct {
	emitter stats.Sender
	log     loggy.Logger
	history string
}

func newHistoryHandler(emitter stats.Sender, history string) http.Handler {
	return &historyHandler{
		emitter: emitter,
		log:     loggy.New("history-handler"),
		history: history,
	}
}

func (h *historyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(h.history))

	h.emitter.Count("ui-history-ok", 1)
}
