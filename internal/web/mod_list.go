package web

import (
	"net/http"
	"strings"

	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/modprox-proxy/internal/modules/store"
	"github.com/pkg/errors"
)

type moduleList struct {
	index store.Index
	log   loggy.Logger
}

func modList(index store.Index) http.Handler {
	return &moduleList{
		index: index,
		log:   loggy.New("mod-list"),
	}
}

func (h *moduleList) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	module, err := moduleFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.log.Infof("serving request for list module: %s", module)

	listing, err := h.index.List(module)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body := formatList(listing)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))

	return
}

func formatList(list []string) string {
	var sb strings.Builder
	for _, version := range list {
		sb.WriteString(version)
		sb.WriteString("\n")
	}
	return sb.String()
}

// GET baseURL/module/@v/list fetches a list of all known versions, one per line.

func moduleFromPath(p string) (string, error) {
	vIdx := strings.Index(p, "@v")
	if vIdx <= 0 {
		return "", errors.Errorf("malformed path: %q")
	}
	return p[0:vIdx], nil
}
