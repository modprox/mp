package web

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_demangle(t *testing.T) {
	try := func(input, exp string) {
		output := demangle(input)
		require.Equal(t, exp, output)
	}

	try("", "")
	try("foo/bar", "foo/bar")
	try("github.com/!burnt!sushi/toml", "github.com/BurntSushi/toml")
	try("foo/bar!", "foo/bar!")
	try("foo!A/bar!B", "foo!A/bar!B")
}
