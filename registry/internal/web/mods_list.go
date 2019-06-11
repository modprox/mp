package web

import (
	"html/template"
	"net/http"
	"sort"

	"oss.indeed.com/go/modprox/pkg/coordinates"
	"oss.indeed.com/go/modprox/pkg/loggy"
	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/registry/internal/data"
	"oss.indeed.com/go/modprox/registry/static"
)

type modsListPage struct {
	Mods map[string][]string // pkg => []version
}

type modsListHandler struct {
	html    *template.Template
	store   data.Store
	emitter stats.Sender
	log     loggy.Logger
}

func newModsListHandler(store data.Store, emitter stats.Sender) http.Handler {
	html := static.MustParseTemplates(
		"static/html/layout.html",
		"static/html/navbar.html",
		"static/html/mods_list.html",
	)

	return &modsListHandler{
		html:    html,
		store:   store,
		emitter: emitter,
		log:     loggy.New("list-modules-handler"),
	}
}

func (h *modsListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code, page, err := h.get(r)
	if err != nil {
		h.log.Errorf("failed to serve modules list page")
		http.Error(w, err.Error(), code)
		h.emitter.Count("ui-list-mods-error", 1)
		return
	}

	if err := h.html.Execute(w, page); err != nil {
		h.log.Errorf("failed to execute modules list page")
		return
	}

	h.emitter.Count("ui-list-mods-ok", 1)
}

func (h *modsListHandler) get(r *http.Request) (int, *modsListPage, error) {
	mods, err := h.store.ListModules()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	tree := treeOfMods(mods)
	page := &modsListPage{
		Mods: tree,
	}

	return http.StatusOK, page, nil
}

func treeOfMods(mods []coordinates.SerialModule) map[string][]string {
	m := make(map[string][]string)
	for _, mod := range mods {
		m[mod.Source] = append(m[mod.Source], mod.Version)
	}

	for _, mod := range mods {
		sort.Strings(m[mod.Source])
	}

	return m
}
