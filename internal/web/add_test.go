package web

import (
	"errors"
	"testing"

	"github.com/modprox/libmodprox/repository"

	"github.com/stretchr/testify/require"
)

func Test_modLineRe(t *testing.T) {
	try := func(input string, exp []string) {
		groups := modLineRe.FindStringSubmatch(input)
		require.Equal(t, exp, groups)
	}

	try( // malformed
		"bad",
		[]string(nil),
	)
	try( // normal
		"github.com/foo/bar v2.0.0",
		[]string{"github.com/foo/bar v2.0.0", "github.com/foo/bar", "v2.0.0", ""},
	)
	try( // with timestamp and hash
		"github.com/foo/bar v0.0.0-20180111040409-fbec762f837d",
		[]string{
			"github.com/foo/bar v0.0.0-20180111040409-fbec762f837d",
			"github.com/foo/bar",
			"v0.0.0-20180111040409-fbec762f837d",
			"-20180111040409-fbec762f837d", // subgroup is flat
		},
	)
	try( // prefix space
		"    github.com/foo/bar v2.0.0",
		[]string{"github.com/foo/bar v2.0.0", "github.com/foo/bar", "v2.0.0", ""},
	)
	try( // comment
		"    github.com/foo/bar v2.0.0 // indirect",
		[]string{"github.com/foo/bar v2.0.0", "github.com/foo/bar", "v2.0.0", ""},
	)
}

func Test_parseLine(t *testing.T) {
	try := func(input string, exp Parsed) {
		parsed := parseLine(input)
		require.Equal(t, exp, parsed)
	}

	try( // malformed
		"github.com/foo/bar",
		Parsed{
			Text:   "github.com/foo/bar",
			Module: repository.ModInfo{},
			Err:    errors.New("malformed module line"),
		},
	)

	try( // normal
		"github.com/foo/bar v2.0.0",
		Parsed{
			Text: "github.com/foo/bar v2.0.0",
			Module: repository.ModInfo{
				Source:  "github.com/foo/bar",
				Version: "v2.0.0",
			},
			Err: nil,
		},
	)

	try( // with timestamp and hash
		"github.com/foo/bar v0.0.0-20180111040409-fbec762f837d",
		Parsed{
			Text: "github.com/foo/bar v0.0.0-20180111040409-fbec762f837d",
			Module: repository.ModInfo{
				Source:  "github.com/foo/bar",
				Version: "v0.0.0-20180111040409-fbec762f837d",
			},
			Err: nil,
		},
	)
}
