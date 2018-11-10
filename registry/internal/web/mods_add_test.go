package web

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
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

	try( // with /go.mod annotation
		"github.com/boltdb/bolt v1.3.1/go.mod h1:clJnj/oiGkjum5o1McbSZDSLxVThjynRyGBgiAx27Ps=",
		Parsed{
			Text: "github.com/boltdb/bolt v1.3.1/go.mod h1:clJnj/oiGkjum5o1McbSZDSLxVThjynRyGBgiAx27Ps=",
			Module: coordinates.Module{
				Source:  "github.com/boltdb/bolt",
				Version: "v1.3.1",
			},
			Err: nil,
		},
	)
}

func Test_parseLines_sumFile(t *testing.T) {
	inputLines := linesOf(t, goSumFile)
	expLines := asMap(linesOf(t, goSumFileExp))

	parsed := parseLines(inputLines)

	for _, m := range parsed {
		s := fmt.Sprintf("%s %s", m.Module.Source, m.Module.Version)
		_, exists := expLines[s]
		require.True(t, exists)
	}
}

func linesOf(t *testing.T, text string) []string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	require.NoError(t, scanner.Err())
	return lines
}

func asMap(lines []string) map[string]struct{} {
	m := make(map[string]struct{}, len(lines))
	for _, line := range lines {
		m[line] = struct{}{}
	}
	return m
}

const goSumFile = `
github.com/boltdb/bolt v1.3.1/go.mod h1:clJnj/oiGkjum5o1McbSZDSLxVThjynRyGBgiAx27Ps=
github.com/cactus/go-statsd-client v3.1.1+incompatible/go.mod h1:cMRcwZDklk7hXp+Law83urTHUiHMzCev/r4JMYr/zU0=
github.com/davecgh/go-spew v1.1.1 h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=
github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/go-sql-driver/mysql v1.4.0/go.mod h1:zAC/RDZ24gD3HViQzih4MyKcchzm+sOG5ZlKdlhCg5w=
github.com/gorilla/context v1.1.1/go.mod h1:kBGZzfjB9CEq2AlWe17Uuf7NDRt0dE0s8S51q0aT7Yg=
github.com/gorilla/csrf v1.5.1/go.mod h1:HTDW7xFOO1aHddQUmghe9/2zTvg7AYCnRCs7MxTGu/0=
github.com/gorilla/mux v1.6.2/go.mod h1:1lud6UwP+6orDFRuTfBEV8e9/aOM/c4fVVCaMa2zaAs=
github.com/gorilla/securecookie v1.1.1/go.mod h1:ra0sb63/xPlUeL+yeDciTfxMRAA+MP+HVt/4epWDjd4=
github.com/jinzhu/copier v0.0.0-20180308034124-7e38e58719c3/go.mod h1:yL958EeXv8Ylng6IfnvG4oflryUi3vgA3xPs9hmII1s=
github.com/lib/pq v1.0.0/go.mod h1:5WUZQaWbwv1U+lTReE5YruASi9Al49XbQIvNi/34Woo=
github.com/modprox/mp v0.0.3 h1:0bVN3YPGWiRYm5gOdQK8slj1Hq4XBf8wctKBu1Mtqlw=
github.com/modprox/mp v0.0.3/go.mod h1:OcYAxmKsY4565H72cVdkhZSh9jmAK1tjlrYAmfrWKFs=
github.com/modprox/mp v0.0.4 h1:WJn8n0ANbfZG5JYpIbcO3OXuUdMSAWPrjOvJG9alDEw=
github.com/modprox/mp v0.0.4/go.mod h1:OcYAxmKsY4565H72cVdkhZSh9jmAK1tjlrYAmfrWKFs=
github.com/modprox/mp v0.0.5 h1:Pb+x8okiUV1Gzp8DDgAEbYUDlaugR/DB0YzUt728ewI=
github.com/modprox/mp v0.0.5/go.mod h1:OcYAxmKsY4565H72cVdkhZSh9jmAK1tjlrYAmfrWKFs=
github.com/pkg/errors v0.8.0 h1:WdK/asTD0HN+q6hsWO3/vpuAkAr+tw6aNJNDFFf0+qw=
github.com/pkg/errors v0.8.0/go.mod h1:bwawxfHBFNV+L2hUp1rHADufV3IMtnDRdf1r5NINEl0=
github.com/pmezard/go-difflib v1.0.0 h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=
github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
github.com/shoenig/atomicfs v0.1.1/go.mod h1:aw1MEIbywlavaTMRvnai/OLCh2dV5XWNcNBYf6iLDno=
github.com/shoenig/petrify/v4 v4.0.2/go.mod h1:xTXgxRisT/LPHgtw0yWpLdVJPyocAlhAPwTKktpm6f4=
github.com/shoenig/toolkit v1.0.0 h1:bevOHX/3xqlV3AGTGkFSYu1a+v8bWMJAZ7kUEj1f7d4=
github.com/shoenig/toolkit v1.0.0/go.mod h1:AzSCIBam5p35X6rgoLpLG/PDQPC6sMUr6nPz8zHWDNk=
github.com/stretchr/objx v0.1.1 h1:2vfRuCMp5sSVIDSqO8oNnWJq7mPa6KVP3iPIwFBuy8A=
github.com/stretchr/objx v0.1.1/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
github.com/stretchr/testify v1.2.2 h1:bSDNvY7ZPG5RlJ8otE/7V6gMiyenm9RtJ7IUVIAoJ1w=
github.com/stretchr/testify v1.2.2/go.mod h1:a8OnRcib4nhh0OaRAV+Yts87kKdq0PP7pXfy6kDkUVs=
golang.org/x/sys v0.0.0-20180909124046-d0be0721c37e/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
google.golang.org/appengine v1.1.0/go.mod h1:EbEs0AVv82hx2wNQdGPgUI5lhzA/G0D9YwlJXL52JkM=
`

const goSumFileExp = `
github.com/boltdb/bolt v1.3.1 
github.com/cactus/go-statsd-client v3.1.1+incompatible 
github.com/davecgh/go-spew v1.1.1
github.com/go-sql-driver/mysql v1.4.0 
github.com/gorilla/context v1.1.1 
github.com/gorilla/csrf v1.5.1 
github.com/gorilla/mux v1.6.2 
github.com/gorilla/securecookie v1.1.1 
github.com/jinzhu/copier v0.0.0-20180308034124-7e38e58719c3 
github.com/lib/pq v1.0.0 
github.com/modprox/mp v0.0.3
github.com/modprox/mp v0.0.4
github.com/modprox/mp v0.0.5
github.com/pkg/errors v0.8.0
github.com/pmezard/go-difflib v1.0.0
github.com/shoenig/atomicfs v0.1.1 
github.com/shoenig/petrify/v4 v4.0.2 
github.com/shoenig/toolkit v1.0.0
github.com/stretchr/objx v0.1.1
github.com/stretchr/testify v1.2.2
golang.org/x/sys v0.0.0-20180909124046-d0be0721c37e 
google.golang.org/appengine v1.1.0 `
