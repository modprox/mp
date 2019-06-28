package upstream

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func dummyRequest() *Request {
	return &Request{
		Transport:     "https",
		Domain:        "code.example.com",
		Namespace:     []string{"go", "foo"},
		Version:       "v0.0.1",
		Path:          "/a/b/c.zip",
		GoGetRedirect: true,
		Headers:       map[string]string{"X-Something": "abc123"},
	}
}

func Test_Request_String(t *testing.T) {
	request := dummyRequest()

	s := request.String()
	require.Equal(t,
		`["https" "code.example.com" [go foo] "v0.0.1" "/a/b/c.zip" true]`,
		s,
	)
}

func Test_Request_URI(t *testing.T) {
	request := dummyRequest()

	uri := request.URI()
	require.Equal(t,
		`https://code.example.com/a/b/c.zip`,
		uri,
	)
}

func Test_Request_Equals_yes(t *testing.T) {
	r1 := dummyRequest()
	r2 := dummyRequest()

	require.True(t, r1.Equals(r2))
	require.True(t, r2.Equals(r1))
}

func Test_Request_Equals_no_transport(t *testing.T) {
	r1 := dummyRequest()
	r2 := dummyRequest()

	r2.Transport = "http"
	require.False(t, r1.Equals(r2))
	require.False(t, r2.Equals(r1))
}

func Test_Request_Equals_no_domain(t *testing.T) {
	r1 := dummyRequest()
	r2 := dummyRequest()

	r2.Domain = "src.example.com"
	require.False(t, r1.Equals(r2))
	require.False(t, r2.Equals(r1))
}

func Test_Request_Equals_no_namespace(t *testing.T) {
	r1 := dummyRequest()
	r2 := dummyRequest()

	r2.Namespace = []string{"x", "y"}
	require.False(t, r1.Equals(r2))
	require.False(t, r2.Equals(r1))
}

func Test_Request_Equals_no_version(t *testing.T) {
	r1 := dummyRequest()
	r2 := dummyRequest()

	r2.Version = "v2.2.2"
	require.False(t, r1.Equals(r2))
	require.False(t, r2.Equals(r1))
}

func Test_Request_Equals_no_path(t *testing.T) {
	r1 := dummyRequest()
	r2 := dummyRequest()

	r2.Path = "/x/y.zip"
	require.False(t, r1.Equals(r2))
	require.False(t, r2.Equals(r1))
}

func Test_Request_Equals_no_goget(t *testing.T) {
	r1 := dummyRequest()
	r2 := dummyRequest()

	r2.GoGetRedirect = false
	require.False(t, r1.Equals(r2))
	require.False(t, r2.Equals(r1))
}

func Test_Request_Equals_no_headers(t *testing.T) {
	r1 := dummyRequest()
	r2 := dummyRequest()

	r2.Headers = map[string]string{"foo": "bar"}
	require.False(t, r1.Equals(r2))
	require.False(t, r2.Equals(r1))
}
