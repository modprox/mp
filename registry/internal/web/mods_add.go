package web

import (
	"errors"
	"html/template"
	"net/http"
	"strings"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/gorilla/csrf"

	"github.com/modprox/mp/pkg/coordinates"
	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/pkg/repository"
	"github.com/modprox/mp/registry/internal/data"
	"github.com/modprox/mp/registry/static"
)

type newPage struct {
	Mods []Parsed
	CSRF template.HTML
}

type newHandler struct {
	html    *template.Template
	store   data.Store
	statter statsd.Statter
	log     loggy.Logger
}

func newAddHandler(store data.Store, statter statsd.Statter) http.Handler {
	html := static.MustParseTemplates(
		"static/html/layout.html",
		"static/html/navbar.html",
		"static/html/mods_add.html",
	)

	return &newHandler{
		html:    html,
		store:   store,
		statter: statter,
		log:     loggy.New("add-modules-handler"),
	}
}

func (h *newHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		code int
		err  error
		page *newPage
	)

	switch r.Method {
	case http.MethodGet:
		code, page, err = h.get(r)
	case http.MethodPost:
		code, page, err = h.post(r)
	}

	if err != nil {
		h.log.Errorf("failed to serve add-module page: %v", err)
		http.Error(w, err.Error(), code)
		_ = h.statter.Inc("ui-add-mod-error", 1, 1)
		return
	}

	if err := h.html.Execute(w, page); err != nil {
		h.log.Errorf("failed to execute add-module page: %v", err)
		return
	}

	_ = h.statter.Inc("ui-add-mod-ok", 1, 1)
}

func (h *newHandler) get(r *http.Request) (int, *newPage, error) {
	return http.StatusOK, &newPage{
		Mods: nil,
		CSRF: csrf.TemplateField(r),
	}, nil
}

func (h *newHandler) post(r *http.Request) (int, *newPage, error) {
	mods, err := h.parseTextArea(r)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	modulesAdded, err := h.storeNewMods(mods)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	h.log.Infof("added %d new modules", modulesAdded)

	return http.StatusOK, &newPage{
		Mods: mods,
		CSRF: csrf.TemplateField(r),
	}, nil
}

func (h *newHandler) storeNewMods(mods []Parsed) (int, error) {
	ableToAdd := make([]coordinates.Module, 0, len(mods))
	for _, parsed := range mods {
		if parsed.Err == nil {
			ableToAdd = append(ableToAdd, parsed.Module)
		}
	}

	for _, able := range ableToAdd {
		h.log.Tracef("[web] adding to registry: %s@%s", able.Source, able.Version)
	}

	return h.store.InsertModules(ableToAdd)
}

type Parsed struct {
	Text   string
	Module coordinates.Module
	Err    error
}

func (h *newHandler) parseTextArea(r *http.Request) ([]Parsed, error) {
	// get the text from form and use a scanner to get each line
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	text := r.PostForm.Get("modules-input")

	// parse the text field into lines then module + tag
	lines := linesOfText(text)
	if len(lines) == 0 {
		return nil, errors.New("no modules listed")
	}
	results := h.parseLines(lines)

	return results, nil
}

func (h *newHandler) parseLines(lines []string) []Parsed {
	results := make([]Parsed, 0, len(lines))
	for _, line := range lines {
		if !h.skipLine(line) {
			result := h.parseLine(line)
			results = append(results, result)
		}
	}
	return results
}

func (h *newHandler) skipLine(line string) bool {
	if strings.HasPrefix(line, "module ") {
		return true
	}
	if strings.Contains(line, "(") {
		return true
	}
	if strings.Contains(line, ")") {
		return true
	}
	return false
}

func (h *newHandler) parseLine(line string) Parsed {
	mod, err := repository.Parse(line)
	return Parsed{
		Text:   line,
		Module: mod,
		Err:    err,
	}
}
