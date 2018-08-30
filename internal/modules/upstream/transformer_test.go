package upstream

import (
	"testing"

	"github.com/modprox/libmodprox/repository"
	"github.com/stretchr/testify/require"
)

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
		Path:      "shoenig/toolkit",
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
		Path:      "a/b/c",
		Version:   "v1.0.1",
	}

	rt := NewRedirectTransform("mycompany", "code.mycompany.net")

	transformed := rt.Modify(request)
	require.Equal(t, &Request{
		Transport: "https",
		Domain:    "code.mycompany.net",
		Path:      "a/b/c",
		Version:   "v1.0.1",
	}, transformed)
}
