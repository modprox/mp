package zips

import (
	"testing"

	"github.com/stretchr/testify/require"
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
