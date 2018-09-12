package web

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/modprox/mp/proxy/config"

	"github.com/modprox/mp/pkg/clients/payloads"
	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/registry/internal/data"
	"github.com/modprox/mp/registry/static"
)

// todo: belongs on some other page
// type linkable struct {
// Module coordinates.Module
// WebURL string
// TagURL string
//}

type ProxyState struct {
	Heartbeat      payloads.Heartbeat
	Configuration  payloads.Configuration
	TransformsText string
}

type homePage struct {
	// Modules []linkable // todo: somewhere else

	Proxies []ProxyState
}

type homeHandler struct {
	html  *template.Template
	store data.Store
	log   loggy.Logger
}

func newHomeHandler(store data.Store) http.Handler {
	html := static.MustParseTemplates(
		"static/html/layout.html",
		"static/html/navbar.html",
		"static/html/home.html",
	)
	return &homeHandler{
		html:  html,
		store: store,
		log:   loggy.New("home-handler"),
	}
}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Tracef("serving up the homepage, path: %s", r.URL.Path)

	configs, err := h.store.ListStartConfigs()
	if err != nil {
		http.Error(w, "failed to list proxy configs", http.StatusInternalServerError)
		h.log.Errorf("failed to list proxy configs: %v", err)
		return
	}

	heartbeats, err := h.store.ListHeartbeats()
	if err != nil {
		http.Error(w, "failed to list proxy heartbeats", http.StatusInternalServerError)
		h.log.Errorf("failed to list proxy heartbeats: %v", err)
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

	// find the associated hearbeat, if any and append them to list of proxy states
	// which will then be used in the template

	// todo: move this to some other page
	//modules, err := h.store.ListModules()
	//if err != nil {
	//	http.Error(w, "failed to list sources", http.StatusInternalServerError)
	//	h.log.Tracef("failed to list sources: %v", err)
	//	return
	//}
	//
	//page := homePage{Modules: linkables(modules)}
	//

	// startupConfigs := h.store.
}

func transformsText(t config.Transforms) string {
	bs, err := json.MarshalIndent(t, "", " ")
	// bs, err := json.Marshal(t)
	if err != nil {
		return "{}"
	}
	return string(bs)
}

// todo: move to some other page
//func linkables(modules []coordinates.SerialModule) []linkable {
//	l := make([]linkable, 0, len(modules))
//	for _, module := range modules {
//		webURL, tagURL := urlInfo(module.Module)
//		l = append(l, linkable{
//			Module: module.Module,
//			WebURL: webURL,
//			TagURL: tagURL,
//		})
//	}
//	return l
//}
//
//func urlInfo(module coordinates.Module) (string, string) {
//	if strings.HasPrefix(module.Source, "github") {
//		webURL := fmt.Sprintf("https://%s", module.Source)
//		tagURL := fmt.Sprintf("https://%s/releases/tag/%s", module.Source, module.Version)
//		return webURL, tagURL
//	}
//	return "#", "#"
//}
