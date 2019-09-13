package zips

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"gophers.dev/pkgs/loggy"
	"gophers.dev/pkgs/semantic"
)

func Test_mangle(t *testing.T) {
	try := func(source, exp string) {
		result := mangle(source)
		require.Equal(t, exp, result)
	}

	try("a.com", "a.com")
	try("A.com", "!a.com")
	try("a.COM", "a.!c!o!m")
	try("alpha.com/foo/bar", "alpha.com/foo/bar")
	try("alpha.com/Foo/BAr", "alpha.com/!foo/!b!ar")
	try("github.com/Azure/azure-sdk-for-go", "github.com/!azure/azure-sdk-for-go")
	try("github.com/GoogleCloudPlatform/cloudsql-proxy", "github.com/!google!cloud!platform/cloudsql-proxy")
	try("github.com/Sirupsen/logrus", "github.com/!sirupsen/logrus")
}

func TestProxyClient_List(t *testing.T) {
	httpClient := NewIHTTPClientMock(t)
	defer httpClient.MinimockFinish()

	const responseBody = `v0.1.0
v0.2.0
v0.3.0
v0.3.1
`

	httpClient.DoMock.Set(func(req *http.Request) (rp1 *http.Response, err error) {
		require.Equal(t, "https://proxy.golang.org/github.com/foo/bar/@v/list", req.URL.String())
		return &http.Response{Body: ioutil.NopCloser(strings.NewReader(responseBody)), StatusCode: http.StatusOK}, nil
	})

	subject := &proxyClient{
		httpClient: httpClient,
		baseURL:    "proxy.golang.org",
		protocol:   "https",
		log:        loggy.New(""),
	}

	versions, err := subject.List("github.com/foo/bar")
	require.NoError(t, err)

	require.Equal(t, []semantic.Tag{
		{Major: 0, Minor: 3, Patch: 1},
		{Major: 0, Minor: 3, Patch: 0},
		{Major: 0, Minor: 2, Patch: 0},
		{Major: 0, Minor: 1, Patch: 0},
	}, versions)
}
