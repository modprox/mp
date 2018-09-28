package web

import (
	"errors"
	"html/template"
	"net/http"
	"sort"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/modprox/mp/pkg/coordinates"
	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/registry/internal/data"
	"github.com/modprox/mp/registry/static"
)

type showPage struct {
	Source string
	Mods   []coordinates.SerialModule
}

type showHandler struct {
	html    *template.Template
	store   data.Store
	statter statsd.Statter
	log     loggy.Logger
}

func newShowHandler(store data.Store, statter statsd.Statter) http.Handler {
	html := static.MustParseTemplates(
		"static/html/layout.html",
		"static/html/navbar.html",
		"static/html/mods_show.html",
	)

	return &showHandler{
		html:    html,
		store:   store,
		statter: statter,
		log:     loggy.New("show-module-h"),
	}
}

func (h *showHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// for now this is get only
	code, page, err := h.get(r)
	if err != nil {
		h.log.Errorf("failed to serve show modules page: %v", err)
		http.Error(w, err.Error(), code)
		h.statter.Inc("ui-show-mod-error", 1, 1)
		return
	}

	if err := h.html.Execute(w, page); err != nil {
		h.log.Errorf("failed to execute show modules page: %v", err)
	}

	h.statter.Inc("ui-show-mod-ok", 1, 1)
}

func (h *showHandler) get(r *http.Request) (int, *showPage, error) {
	source, err := h.parseQuery(r)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	mods, err := h.store.ListModulesBySource(source)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	sort.Sort(coordinates.ModsByVersion(mods))

	return http.StatusOK, &showPage{
		Source: source,
		Mods:   mods,
	}, nil
}

func (h *showHandler) parseQuery(r *http.Request) (string, error) {
	values := r.URL.Query()
	m := values.Get("mod")
	if m == "" {
		return "", errors.New("mod query parameter required")
	}
	return m, nil
}
