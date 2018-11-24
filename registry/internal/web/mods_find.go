package web

import (
	"errors"
	"html/template"
	"net/http"
	"regexp"
	"time"

	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/registry/internal/tools/finder"
	"github.com/modprox/mp/registry/static"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/gorilla/csrf"
)

type findPage struct {
	CSRF  template.HTML
	Found []findResult
}

type findHandler struct {
	html    *template.Template
	statter statsd.Statter
	finder  finder.Finder
	log     loggy.Logger
}

func newFindHandler(statter statsd.Statter) http.Handler {
	html := static.MustParseTemplates(
		"static/html/layout.html",
		"static/html/navbar.html",
		"static/html/mods_find.html",
	)

	return &findHandler{
		html:    html,
		statter: statter,
		finder: finder.New(finder.Options{
			Timeout: 1 * time.Minute,
		}),
		log: loggy.New("find-modules-handler"),
	}
}

func (h *findHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		code int
		err  error
		page *findPage
	)

	switch r.Method {
	case http.MethodGet:
		code, page, err = h.get(r)
	case http.MethodPost:
		code, page, err = h.post(r)
	}

	if err != nil {
		h.log.Errorf("failed to serve find-module page: %v", err)
		http.Error(w, err.Error(), code)
		_ = h.statter.Inc("ui-find-mod-error", 1, 1)
		return
	}

	if err := h.html.Execute(w, page); err != nil {
		h.log.Errorf("failed to execute find-module page: %v", err)
		return
	}

	_ = h.statter.Inc("ui-find-mod-ok", 1, 1)
}

func (h *findHandler) get(r *http.Request) (int, *findPage, error) {
	return http.StatusOK, &findPage{
		CSRF: csrf.TemplateField(r),
	}, nil
}

func (h *findHandler) post(r *http.Request) (int, *findPage, error) {
	results, err := h.parseTextArea(r)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	return http.StatusOK, &findPage{
		CSRF:  csrf.TemplateField(r),
		Found: results,
	}, nil
}

type findResult struct {
	Text   string
	Result *finder.Result
	Err    error
}

func (h *findHandler) parseTextArea(r *http.Request) ([]findResult, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	text := r.PostForm.Get("sources-input")

	lines := linesOfText(text)
	if len(lines) == 0 {
		return nil, errors.New("no sources listed")
	}

	results := h.processLines(lines)

	return results, nil
}

func (h *findHandler) processLines(lines []string) []findResult {
	results := make([]findResult, 0, len(lines))
	for _, line := range lines {
		result := h.processLine(line)
		results = append(results, result)
	}
	return results
}

// only github.com things are supported for now
var findableRe = regexp.MustCompile(`github\.com/[[:alnum:]]+/[[:alnum:]]+`)

func (h *findHandler) processLine(line string) findResult {
	if !findableRe.MatchString(line) {
		return findResult{
			Text: line,
			Err:  errors.New("does not match regexp"),
		}
	}

	result, err := h.finder.Find(line)
	if err != nil {
		h.log.Warnf("failed to find result for %s: %v", line, err)
		return findResult{
			Text: line,
			Err:  err,
		}
	}

	return findResult{
		Text:   line,
		Result: result,
	}
}
