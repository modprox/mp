package web

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/libmodprox/repository"
	"github.com/modprox/modprox-registry/internal/data"
	"github.com/modprox/modprox-registry/static"
)

type linkable struct {
	Module repository.ModInfo
	WebURL string
	TagURL string
}

type homePage struct {
	Modules []linkable
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

	modules, err := h.store.ListMods()
	if err != nil {
		http.Error(w, "failed to list sources", http.StatusInternalServerError)
		h.log.Tracef("failed to list sources: %v", err)
		return
	}

	page := homePage{Modules: linkables(modules)}

	if err := h.html.Execute(w, page); err != nil {
		h.log.Errorf("failed to execute homepage template: %v", err)
		return
	}
}

func linkables(modules []repository.ModInfo) []linkable {
	l := make([]linkable, 0, len(modules))
	for _, module := range modules {
		webURL, tagURL := urlInfo(module)
		l = append(l, linkable{
			Module: module,
			WebURL: webURL,
			TagURL: tagURL,
		})
	}
	return l
}

func urlInfo(module repository.ModInfo) (string, string) {
	if strings.HasPrefix(module.Source, "github") {
		webURL := fmt.Sprintf("https://%s", module.Source)
		tagURL := fmt.Sprintf("https://%s/releases/tag/%s", module.Source, module.Version)
		return webURL, tagURL
	}
	return "#", "#"
}
