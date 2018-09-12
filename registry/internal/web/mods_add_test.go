package web

import (
	"errors"
	"testing"

	"github.com/modprox/mp/pkg/coordinates"

	"github.com/stretchr/testify/require"
)

func compareErr(t *testing.T, expErr, gotErr error) {
	if expErr == nil {
		require.Nil(t, gotErr)
	} else {
		require.NotNil(t, gotErr)
		require.Equal(t, expErr.Error(), gotErr.Error())
	}
}

func Test_parseLine(t *testing.T) {
	try := func(input string, exp Parsed) {
		parsed := parseLine(input)
		require.Equal(t, exp.Text, parsed.Text)
		require.Equal(t, exp.Module, parsed.Module)
		compareErr(t, exp.Err, parsed.Err)
	}

	try( // malformed
		"github.com/foo/bar",
		Parsed{
			Text:   "github.com/foo/bar",
			Module: coordinates.Module{},
			Err:    errors.New(`malformed module line: "github.com/foo/bar"`),
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

	try( // with @ notation
		"github.com/kr/pty@v1.1.1",
		Parsed{
			Text: "github.com/kr/pty@v1.1.1",
			Module: coordinates.Module{
				Source:  "github.com/kr/pty",
				Version: "v1.1.1",
			},
			Err: nil,
		},
	)
}
