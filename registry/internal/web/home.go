package web

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/modprox/mp/pkg/clients/payloads"
	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/pkg/metrics/stats"
	"github.com/modprox/mp/proxy/config"
	"github.com/modprox/mp/registry/internal/data"
	"github.com/modprox/mp/registry/static"
)

type ProxyState struct {
	Heartbeat      payloads.Heartbeat
	Configuration  payloads.Configuration
	TransformsText string
}

type homePage struct {
	Proxies []ProxyState
}

type homeHandler struct {
	html    *template.Template
	store   data.Store
	emitter stats.Sender
	log     loggy.Logger
}

func newHomeHandler(store data.Store, emitter stats.Sender) http.Handler {
	html := static.MustParseTemplates(
		"static/html/layout.html",
		"static/html/navbar.html",
		"static/html/home.html",
	)
	return &homeHandler{
		html:    html,
		store:   store,
		emitter: emitter,
		log:     loggy.New("home-handler"),
	}
}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	configs, err := h.store.ListStartConfigs()
	if err != nil {
		http.Error(w, "failed to list proxy configs", http.StatusInternalServerError)
		h.log.Errorf("failed to list proxy configs: %v", err)
		h.emitter.Count("ui-home-error", 1)
		return
	}

	heartbeats, err := h.store.ListHeartbeats()
	if err != nil {
		http.Error(w, "failed to list proxy heartbeats", http.StatusInternalServerError)
		h.log.Errorf("failed to list proxy heartbeats: %v", err)
		h.emitter.Count("ui-home-error", 1)
		return
	}

	var proxyStates []ProxyState
	for _, c := range configs { // could be more efficient
		state := ProxyState{
			Configuration:  c,
			TransformsText: transformsText(c.Transforms),
		}
		for _, h := range heartbeats {
			if c.Self == h.Self {
				state.Heartbeat = h
				break
			}
		}
		proxyStates = append(proxyStates, state)
	}

	page := homePage{
		Proxies: proxyStates,
	}

	if err := h.html.Execute(w, page); err != nil {
		h.log.Errorf("failed to execute homepage template: %v", err)
		return
	}

	h.emitter.Count("ui-home-ok", 1)
}

func transformsText(t config.Transforms) string {
	bs, err := json.MarshalIndent(t, "", " ")
	if err != nil {
		return "{}"
	}
	return string(bs)
}
