package finder

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"gophers.dev/pkgs/semantic"

	"oss.indeed.com/go/modprox/pkg/clients/zips"
)

func Test_Compatible(t *testing.T) {
	try := func(input string, exp bool) {
		result := Compatible(input)
		require.Equal(t, exp, result)
	}

	// github OK
	try("github.com/foo/bar", true)
	try("github.com/foo/bar/baz", true)
	try("github.com/sean-/seed", true)

	// github NOT OK
	try("github.com/foo", false)
	try("github.com", false)
	try("github", false)

	// nothing else supported
	try("golang.org/x/y", false)
	try("", false)
}

func Test_finder_Find(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.String(), "/tags") {
				_, err := w.Write([]byte(tags))
				require.NoError(t, err)
			} else {
				_, err := w.Write([]byte(head))
				require.NoError(t, err)
			}
		}),
	)
	defer ts.Close()

	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	const source = "github.com/octocat/Hello-World"

	proxyClient := zips.NewProxyClientMock(t)
	defer proxyClient.MinimockFinish()
	proxyClient.ListMock.Expect(source).Return([]semantic.Tag{}, nil)

	f := New(Options{
		Timeout: 1 * time.Second,
		Versions: map[string]Versions{
			"github.com": Github(ts.URL, client, proxyClient),
		},
	})

	result, err := f.Find(source)
	require.NoError(t, err)

	t.Logf("result %#v", result)
}

const tags = `
[
  {
    "name": "v0.1",
    "commit": {
      "sha": "c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc",
      "url": "https://api.github.com/repos/octocat/Hello-World/commits/c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc"
    },
    "zipball_url": "https://github.com/octocat/Hello-World/zipball/v0.1",
    "tarball_url": "https://github.com/octocat/Hello-World/tarball/v0.1"
  }
]`

const head = `
{
  "sha": "eaae6f7b3e4bb6b3337c1181557e1d44c48235fe",
  "node_id": "MDY6Q29tbWl0MTQxOTE5Mzk5OmVhYWU2ZjdiM2U0YmI2YjMzMzdjMTE4MTU1N2UxZDQ0YzQ4MjM1ZmU=",
  "commit": {
    "author": {
      "name": "Seth Hoenig",
      "email": "hoenig@indeed.com",
      "date": "2018-11-16T20:32:56Z"
    }
   }
}`
