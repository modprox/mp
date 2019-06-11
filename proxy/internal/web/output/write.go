package output

import (
	"encoding/json"
	"net/http"
	"strconv"

	"oss.indeed.com/go/modprox/pkg/repository"
)

const (
	headerContentType             = "Content-Type"
	headerContentDescription      = "Content-Description"
	headerContentTransferEncoding = "Content-Transfer-Encoding"
	headerContentLength           = "Content-Length"
)

const (
	JSON         = "application/json"
	Text         = "text/plain"
	Zip          = "application/zip"
	OctetStream  = "application/octet-stream"
	FileTransfer = "File Transfer"
	Binary       = "binary"
)

func Write(w http.ResponseWriter, mime, content string) {
	w.Header().Set(headerContentType, mime)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(content))
}

func WriteJSON(w http.ResponseWriter, i interface{}) {
	w.Header().Set(headerContentType, JSON)
	if err := json.NewEncoder(w).Encode(i); err != nil {
		panic("failure to encode json response: " + err.Error())
	}
}

func WriteZip(w http.ResponseWriter, blob repository.Blob) {
	w.Header().Set(headerContentType, Zip)
	w.Header().Add(headerContentType, OctetStream)
	w.Header().Set(headerContentDescription, FileTransfer)
	w.Header().Set(headerContentTransferEncoding, Binary)
	w.Header().Set(headerContentLength, strconv.Itoa(len(blob)))

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(blob)
}
