package web

import (
	"net/http"
	"strings"

	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/modprox-proxy/internal/modules/store"
	"github.com/modprox/modprox-proxy/internal/web/output"
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

	listing, err := h.index.Versions(module)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output.Write(w, output.Text, formatList(listing))
}

func formatList(list []string) string {
	var sb strings.Builder
	for _, version := range list {
		sb.WriteString(version)
		sb.WriteString("\n")
	}
	return sb.String()
}
