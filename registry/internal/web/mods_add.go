package web

import (
	"bufio"
	"errors"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/csrf"

	"github.com/modprox/mp/pkg/coordinates"
	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/pkg/repository"
	"github.com/modprox/mp/registry/internal/data"
	"github.com/modprox/mp/registry/static"
)

type newPage struct {
	Mods      []Parsed
	CSRFField template.HTML
}

type newHandler struct {
	html  *template.Template
	store data.Store
	log   loggy.Logger
}

func newAddHandler(store data.Store) http.Handler {
	html := static.MustParseTemplates(
		"static/html/layout.html",
		"static/html/navbar.html",
		"static/html/mods_add.html",
	)

	return &newHandler{
		html:  html,
		store: store,
		log:   loggy.New("add-modules-handler"),
	}
}

func (h *newHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Tracef("loaded page %v", r.Method)

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
		return
	}

	if err := h.html.Execute(w, page); err != nil {
		h.log.Errorf("failed to execute add-module page: %v", err)
	}
}

func (h *newHandler) get(r *http.Request) (int, *newPage, error) {
	return http.StatusOK, &newPage{
		Mods:      nil,
		CSRFField: csrf.TemplateField(r),
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
		Mods:      mods,
		CSRFField: csrf.TemplateField(r),
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
	results := parseLines(lines)

	return results, nil
}

func parseLines(lines []string) []Parsed {
	results := make([]Parsed, 0, len(lines))
	for _, line := range lines {
		if strings.Contains(line, "/go.mod ") {
			// when copying from go.sum, every line
			// appears twice, once with this key so
			// just get rid of this one
			continue
		}
		result := parseLine(line)
		results = append(results, result)
	}
	return results
}

func parseLine(line string) Parsed {
	mod, err := repository.Parse(line)
	return Parsed{
		Text:   line,
		Module: mod,
		Err:    err,
	}
}

func linesOfText(text string) []string {
	lines := make([]string, 0, 1)
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, scanner.Text())
		}
	}
	return lines
}
