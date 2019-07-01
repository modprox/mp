package web

import (
	"errors"
	"html/template"
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/csrf"

	"go.gophers.dev/pkgs/loggy"

	"oss.indeed.com/go/modprox/pkg/coordinates"
	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/registry/internal/data"
	"oss.indeed.com/go/modprox/registry/static"
)

type showPage struct {
	CSRF   template.HTML
	Source string
	Mods   []coordinates.SerialModule
}

type showHandler struct {
	html    *template.Template
	store   data.Store
	emitter stats.Sender
	log     loggy.Logger
}

func newShowHandler(store data.Store, emitter stats.Sender) http.Handler {
	html := static.MustParseTemplates(
		"static/html/layout.html",
		"static/html/navbar.html",
		"static/html/mods_show.html",
	)

	return &showHandler{
		html:    html,
		store:   store,
		emitter: emitter,
		log:     loggy.New("show-module-h"),
	}
}

func (h *showHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		code int
		page *showPage
		err  error
	)

	switch r.Method {
	case http.MethodGet:
		code, page, err = h.get(r)
	case http.MethodPost:
		code, page, err = h.post(r)
	}

	if err != nil {
		h.log.Errorf("failed to serve show modules page: %v", err)
		http.Error(w, err.Error(), code)
		h.emitter.Count("ui-show-mod-error", 1)
		return
	}

	if err := h.html.Execute(w, page); err != nil {
		h.log.Errorf("failed to execute show modules page: %v", err)
	}

	h.emitter.Count("ui-show-mod-ok", 1)
}

func (h *showHandler) get(r *http.Request) (int, *showPage, error) {
	return h.load(r)
}

func (h *showHandler) post(r *http.Request) (int, *showPage, error) {
	id, err := h.parseModToDelete(r)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	h.log.Infof("will delete module of id: %d", id)
	if err := h.store.DeleteModuleByID(id); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// after deletion just load the show page for that package again
	return h.load(r)
}

// both get and post will load the mod show page
// which can be rendered with this load function
func (h *showHandler) load(r *http.Request) (int, *showPage, error) {
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
		CSRF:   csrf.TemplateField(r),
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

func (h *showHandler) parseModToDelete(r *http.Request) (int, error) {
	idText := r.FormValue("delete-mod-id")
	return strconv.Atoi(idText)
}
