package finder

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_finder_Find(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.String(), "/tags") {
				w.Write([]byte(tags))
			} else {
				w.Write([]byte(head))
			}
		}),
	)
	defer ts.Close()

	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	f := New(Options{
		Timeout: 1 * time.Second,
		Versions: map[string]Versions{
			"github.com": Github(ts.URL, client),
		},
	})

	result, err := f.Find("github.com/octocat/Hello-World")
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
