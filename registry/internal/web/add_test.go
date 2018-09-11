package web

import (
	"errors"
	"testing"

	"github.com/modprox/mp/pkg/coordinates"

	"github.com/stretchr/testify/require"
)

func Test_parseLine(t *testing.T) {
	try := func(input string, exp Parsed) {
		parsed := parseLine(input)
		require.Equal(t, exp, parsed)
	}

	try( // malformed
		"github.com/foo/bar",
		Parsed{
			Text:   "github.com/foo/bar",
			Module: coordinates.Module{},
			Err:    errors.New("malformed module line"),
		},
	)

	try( // normal
		"github.com/foo/bar v2.0.0",
		Parsed{
			Text: "github.com/foo/bar v2.0.0",
			Module: coordinates.Module{
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
			Module: coordinates.Module{
				Source:  "github.com/foo/bar",
				Version: "v0.0.0-20180111040409-fbec762f837d",
			},
			Err: nil,
		},
	)
}
