package zips

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_majorVersion(t *testing.T) {
	tryMajorVersion(t, "v2.0.4", "v2", false)
	tryMajorVersion(t, "v1.0.0", "", false)
	tryMajorVersion(t, "v0.0.1", "", false)
	tryMajorVersion(t, "blah", "", true)
	tryMajorVersion(t, "v123", "", true)
}

func tryMajorVersion(t *testing.T, version, expectedMajor string, expectError bool) {
	m, err := majorVersion(version)
	if expectError {
		require.Error(t, err)
	} else {
		require.NoError(t, err)
		require.Equal(t, expectedMajor, m)
	}
}

func Test_ModulePath(t *testing.T) {
	gomod := `module github.com/modprox/mp

require (
	github.com/googleapis/gax-go/v2 v2.0.4
	google.golang.org/grpc 1.19.0
)
`
	expected := "github.com/modprox/mp"

	s := ModulePath([]byte(gomod))

	require.Equal(t, expected, s)
}

func Test_ModulePath_none(t *testing.T) {
	gomod := `// I absent-mindedly commented out the module line
//module github.com/modprox/mp`

	s := ModulePath([]byte(gomod))
	require.Equal(t, "", s)
}

func Test_moduleOf(t *testing.T) {
	goModPath := map[string]string{
		"github.com/billsmith/module1-1.1.4": "github.com/billsmith/module1",
		"github.com/billsmith/module1/v2": "github.com/billsmith/module1/v2",
		"github.com/billsmith/module1/v3": "github.com/billsmith/module1/v3",
	}

	require.Equal(t, "github.com/billsmith/module1", moduleOf(goModPath, "github.com/billsmith/module1-1.1.4/main.go"))
}

func Test_moduleOf_v2(t *testing.T) {
	goModPath := map[string]string{
		"github.com/billsmith/module1-2.0.9": "github.com/billsmith/module1",
		"github.com/billsmith/module1-2.0.9/v2": "github.com/billsmith/module1/v2",
		"github.com/billsmith/module1-2.0.9/v3": "github.com/billsmith/module1/v3",
	}

	require.Equal(t, "github.com/billsmith/module1/v2", moduleOf(goModPath, "github.com/billsmith/module1-2.0.9/v2/main.go"))
}

func Test_moduleOf_v4(t *testing.T) {
	goModPath := map[string]string{
		"github.com/billsmith/module1/v2": "github.com/billsmith/module1/v2",
		"github.com/billsmith/module1/v3": "github.com/billsmith/module1/v3",
	}

	require.Equal(t, "", moduleOf(goModPath, "github.com/billsmith/module1/v4/main.go"))
}