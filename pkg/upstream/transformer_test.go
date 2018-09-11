package upstream

import (
	"strings"
	"testing"

	"github.com/modprox/mp/pkg/coordinates"

	"github.com/stretchr/testify/require"
)

func ns(path string) Namespace {
	return strings.Split(path, "/")
}

func Test_NewRequest(t *testing.T) {
	mod := coordinates.Module{
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

func Test_NewRequest_no_path_is_ok(t *testing.T) {
	// An example pulled from real life, go.opencensus.io is
	// itself pointed at a repository using go-get meta.
	// It has no path.
	mod := coordinates.Module{
		Source:  "go.opencensus.io",
		Version: "v0.15.0",
	}

	request, err := NewRequest(mod)
	require.NoError(t, err)
	require.Equal(t, &Request{
		Transport: "https",
		Domain:    "go.opencensus.io",
		Namespace: nil,
		Version:   "v0.15.0",
	}, request)
}

func Test_StaticRedirectTransform(t *testing.T) {
	request := &Request{
		Transport: "https",
		Domain:    "mycompany",
		Namespace: ns("a/b/c"),
		Version:   "v1.0.1",
	}

	rt := NewStaticRedirectTransform("mycompany", "code.mycompany.net")

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

func Test_addressableVersion(t *testing.T) {
	try := func(input string, exp string) {
		output := addressableVersion(input)
		require.Equal(t, exp, output)
	}

	try("v2.0.0", "v2.0.0")
	try("v0.0.0-20180111040409-fbec762f837d", "fbec762f837d")
	try("v2.3.3+incompatible", "v2.3.3")
}
