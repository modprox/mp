package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_about_ok(t *testing.T) {
	h, mocks := makeRouter(t)
	defer mocks.assertions()

	request, err := http.NewRequest(http.MethodGet, "/configure/about", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	h.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
}

func Test_about_bad_method(t *testing.T) {
	h, mocks := makeRouter(t)
	defer mocks.assertions()

	request, err := http.NewRequest(http.MethodPost, "/configure/about", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	h.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
}
