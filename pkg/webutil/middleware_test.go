package webutil

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testHandler struct{}

func (h testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}

func Test_Chain(t *testing.T) {
	var h testHandler
	var called bool
	testMiddleware := func(h http.Handler) http.Handler {
		called = true
		return h
	}
	_ = Chain(h, testMiddleware)
	assert.True(t, called)
}

func Test_KeyGuard_no_header(t *testing.T) {
	guard := KeyGuard([]string{"abc123"})

	executed := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("secret stuff"))
		executed = true
	})

	protected := Chain(handler, guard)

	request, err := http.NewRequest(http.MethodGet, "/foo", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	protected.ServeHTTP(recorder, request)

	code := recorder.Code
	require.Equal(t, http.StatusForbidden, code)
	require.False(t, executed)
}

func Test_KeyGuard_bad_keys(t *testing.T) {
	guard := KeyGuard([]string{"abc123"})

	executed := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("secret stuff"))
		executed = true
	})

	protected := Chain(handler, guard)

	request, err := http.NewRequest(http.MethodGet, "/foo", nil)
	require.NoError(t, err)
	request.Header.Add(HeaderAPIKey, "foo123")
	request.Header.Add(HeaderAPIKey, "bar123")
	request.Header.Add(HeaderAPIKey, "baz123")

	recorder := httptest.NewRecorder()
	protected.ServeHTTP(recorder, request)

	code := recorder.Code
	require.Equal(t, http.StatusForbidden, code)
	require.False(t, executed)
}

func Test_KeyGuard_no_keys(t *testing.T) {
	guard := KeyGuard([]string{""})

	executed := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("secret stuff"))
		executed = true
	})

	protected := Chain(handler, guard)

	request, err := http.NewRequest(http.MethodGet, "/foo", nil)
	require.NoError(t, err)
	request.Header.Add(HeaderAPIKey, "foo123")
	request.Header.Add(HeaderAPIKey, "bar123")
	request.Header.Add(HeaderAPIKey, "baz123")

	recorder := httptest.NewRecorder()
	protected.ServeHTTP(recorder, request)

	code := recorder.Code
	require.Equal(t, http.StatusForbidden, code)
	require.False(t, executed)
}

func Test_KeyGuard_good_key(t *testing.T) {
	guard := KeyGuard([]string{"abc123"})

	executed := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("secret stuff"))
		executed = true
	})

	protected := Chain(handler, guard)

	request, err := http.NewRequest(http.MethodGet, "/foo", nil)
	require.NoError(t, err)
	request.Header.Add(HeaderAPIKey, "foo123")
	request.Header.Add(HeaderAPIKey, "abc123")
	request.Header.Add(HeaderAPIKey, "baz123")

	recorder := httptest.NewRecorder()
	protected.ServeHTTP(recorder, request)

	code := recorder.Code
	require.Equal(t, http.StatusOK, code)
	require.True(t, executed)
}
