package output

import "net/http"

const (
	headerContentType = "Content-Type"
)

const (
	JSON = "application/json"
	Text = "text/plain"
)

// w.Header().Set("Content-Type", "application/json")

func Write(w http.ResponseWriter, mime, content string) {
	w.Header().Set(headerContentType, mime)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(content))
}
