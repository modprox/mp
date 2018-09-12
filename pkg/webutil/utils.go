package webutil

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Write the given struct as JSON into w.
func WriteJSON(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(i); err != nil {
		log.Println("failed to write json response: " + err.Error())
	}
}

// ParseURL parses tsURL, triggering a failure on t if it is not
// possible to do so.
func ParseURL(t *testing.T, tsURL string) (string, int) {
	tsURL = strings.TrimPrefix(tsURL, "http://")
	tokens := strings.Split(tsURL, ":")
	port, err := strconv.Atoi(tokens[1])
	require.NoError(t, err)
	return tokens[0], port
}
