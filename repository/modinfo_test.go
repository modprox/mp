package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parse(t *testing.T) {
	try := func(input string, expMod ModInfo, expErr bool) {
		output, err := Parse(input)
		require.True(t, expErr == (err != nil), "err was: %v", err)
		require.Equal(t, expMod, output)
	}

	try("github.com/foo/bar v2.0.0", ModInfo{
		Source:  "github.com/foo/bar",
		Version: "v2.0.0",
	}, false)

	try("github.com/tdewolff/parse v2.3.3+incompatible // indirect", ModInfo{
		Source:  "github.com/tdewolff/parse",
		Version: "v2.3.3",
	}, false)

	try("golang.org/x/tools v0.0.0-20180111040409-fbec762f837d", ModInfo{
		Source:  "golang.org/x/tools",
		Version: "v0.0.0-20180111040409-fbec762f837d",
	}, false)

	try("/github.com/cpuguy83/go-md2man/@v/v1.0.6.info", ModInfo{
		Source:  "github.com/cpuguy83/go-md2man",
		Version: "v1.0.6",
	}, false)
}

// http://localhost:9000/gopkg.in/check.v1/@v/v0.0.0-20161208181325-20d25e280405.info
