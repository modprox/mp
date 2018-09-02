package service

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testHandler struct {}
func (h testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}

func Test_chain(t *testing.T) {
	var h testHandler
	var called bool
	testMiddleware := func(h http.Handler) http.Handler {
		called = true
		return h
	}
	_ = chain(h, testMiddleware)
	assert.True(t, called)
}