package upstream

import (
	"strings"
	"testing"

	"github.com/modprox/libmodprox/repository"
	"github.com/stretchr/testify/require"
)

func ns(path string) Namespace {
	return strings.Split(path, "/")
}

func Test_NewRequest(t *testing.T) {
	mod := repository.ModInfo{
		Source:  "github.com/shoenig/toolkit",
		Version: "v1.0.1",
	}

	request, err := NewRequest(mod)
	require.NoError(t, err)
	require.Equal(t, &Request{
		Transport: "https",
		Domain:    "github.com",
		Namespace: ns("shoenig/toolkit"),
		Version:   "v1.0.1",
	}, request)
}

func Test_NewRequest_malformed(t *testing.T) {
	mod := repository.ModInfo{
		Source:  "foobar",
		Version: "v1.0.1",
	}

	_, err := NewRequest(mod)
	require.EqualError(t, err, "source does not contain a path")
}

func Test_RedirectTransform(t *testing.T) {
	request := &Request{
		Transport: "https",
		Domain:    "mycompany",
		Namespace: ns("a/b/c"),
		Version:   "v1.0.1",
	}

	rt := NewRedirectTransform("mycompany", "code.mycompany.net")

	transformed := rt.Modify(request)
	require.Equal(t, &Request{
		Transport: "https",
		Domain:    "code.mycompany.net",
		Namespace: ns("a/b/c"),
		Version:   "v1.0.1",
	}, transformed)
}

func Test_formatPath(t *testing.T) {
	try := func(pathFmt string, ns Namespace, version, exp string) {
		result := formatPath(pathFmt, version, ns)
		require.Equal(t, exp, result)
	}

	// github
	try(
		"ELEM1/ELEM2/archive/VERSION.zip",
		ns("shoenig/toolkit"),
		"v1.0.1",
		"shoenig/toolkit/archive/v1.0.1.zip",
	)

	// gitlab
	try(
		"ELEM1/ELEM2/-/archive/VERSION/ELEM2-VERSION.zip",
		ns("crypo/cryptsetup"),
		"v2.0.1",
		"crypo/cryptsetup/-/archive/v2.0.1/cryptsetup-v2.0.1.zip",
	)
}
