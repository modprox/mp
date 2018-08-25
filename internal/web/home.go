package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/modprox/modprox-registry/internal/repositories"
	"github.com/modprox/modprox-registry/internal/repositories/repository"
	"github.com/modprox/modprox-registry/static"
)

type linkable struct {
	Module repository.Module
	WebURL string
	TagURL string
}

type homepage struct {
	Modules []linkable
}

func newHomepageHandler(store repositories.Store) http.Handler {
	html := static.MustParseTemplates(
		"static/html/layout.html",
		"static/html/navbar.html",
		"static/html/home.html",
	)
	return &homepageHandler{
		html:  html,
		store: store,
	}
}

type homepageHandler struct {
	html  *template.Template
	store repositories.Store
}

func (h *homepageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("[web] serving up the homepage")

	modules, err := h.store.ListSources()
	if err != nil {
		http.Error(w, "failed to list sources", http.StatusInternalServerError)
		log.Println("[web] failed to list sources:", err)
	}

	page := homepage{Modules: linkables(modules)}

	if err := h.html.Execute(w, page); err != nil {
		log.Panicf("failed to execute homepage template: %v", err)
	}
}

func linkables(modules []repository.Module) []linkable {
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

func urlInfo(module repository.Module) (string, string) {
	if strings.HasPrefix(module.Source, "github") {
		webURL := fmt.Sprintf("https://%s", module.Source)
		tagURL := fmt.Sprintf("https://%s/releases/tag/%s", module.Source, module.Version)
		return webURL, tagURL
	}
	return "#", "#"
}
