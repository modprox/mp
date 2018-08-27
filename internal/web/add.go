package web

import (
	"bufio"
	"errors"
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/libmodprox/repository"
	"github.com/modprox/modprox-registry/internal/data"
	"github.com/modprox/modprox-registry/static"
)

type newPage struct {
	Mods []Parsed
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
		"static/html/add.html",
	)

	return &newHandler{
		html:  html,
		store: store,
		log:   loggy.New("add-module-handler"),
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
		h.log.Tracef("failed to serve add-module page: %v", err)
		http.Error(w, err.Error(), code)
		return
	}

	if err := h.html.Execute(w, page); err != nil {
		h.log.Tracef("failed to serve add-module page: %v", err)
	}
}

func (h *newHandler) get(r *http.Request) (int, *newPage, error) {
	return http.StatusOK, &newPage{
		Mods: nil,
	}, nil
}

func (h *newHandler) post(r *http.Request) (int, *newPage, error) {
	mods, err := h.parseTextArea(r)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	sourcesAdded, tagsAdded, err := h.storeNewMods(mods)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	h.log.Tracef("added %d tags across %d sources:", tagsAdded, sourcesAdded)

	return http.StatusOK, &newPage{
		Mods: mods,
	}, nil
}

func (h *newHandler) storeNewMods(mods []Parsed) (int, int, error) {
	ableToAdd := make([]repository.ModInfo, 0, len(mods))
	for _, parsed := range mods {
		if parsed.Err == nil {
			ableToAdd = append(ableToAdd, parsed.Module)
		}
	}

	for _, able := range ableToAdd {
		h.log.Tracef("[web] adding to registry: %s@%s", able.Source, able.Version)
	}

	return h.store.AddMod(ableToAdd)
}

type Parsed struct {
	Text   string
	Module repository.ModInfo
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

var (
	// e.g. gihub.com/foo/bar v2.0.0
	modLineRe = regexp.MustCompile(`^([\S.-]+)[\s]+(v[\d]+.[\d]+.[\d]+)$`)
)

func parseLine(line string) Parsed {
	groups := modLineRe.FindStringSubmatch(line)
	if len(groups) != 3 {
		return Parsed{
			Text: line,
			Err:  errors.New("malformed module and tag"),
		}
	}
	return Parsed{
		Text: line,
		Module: repository.ModInfo{
			Source:  strings.TrimRight(groups[1], "/"),
			Version: groups[2],
		},
	}
}

func parseLines(lines []string) []Parsed {
	results := make([]Parsed, 0, len(lines))
	for _, line := range lines {
		result := parseLine(line)
		results = append(results, result)
	}
	return results
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
