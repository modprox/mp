package web

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_history_ok(t *testing.T) {
	h, mocks := makeRouter(t)
	defer mocks.assertions()

	request, err := http.NewRequest(http.MethodGet, "/history", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	h.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)

	bytes, err := ioutil.ReadAll(recorder.Result().Body)
	require.NoError(t, err)
	require.Equal(t, "this is some fake history", string(bytes))
}
