package repository

import (
	"strings"

	"oss.indeed.com/go/modprox/pkg/coordinates"

	"github.com/pkg/errors"
)

var (
// examples
//  mod file style
//    github.com/foo/bar v2.0.0
//    github.com/tdewolff/parse v2.3.3+incompatible // indirect
//    golang.org/x/tools v0.0.0-20180111040409-fbec762f837d
//    gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405
//  sum file style
//    github.com/boltdb/bolt v1.3.1/go.mod h1:clJnj/oiGkjum5o1McbSZDSLxVThjynRyGBgiAx27Ps=
//  proxy request style
//    /github.com/cpuguy83/go-md2man/@v/v1.0.6.info
//  zip style
//    github.com/kr/pty@v1.1.1
)

// Parse will parse s as a module in string form.
func Parse(s string) (coordinates.Module, error) {
	orig := s
	s = strings.Trim(s, "/")
	s = strings.TrimSuffix(s, ".info")
	s = strings.TrimSuffix(s, ".zip")
	s = strings.TrimSuffix(s, ".mod")
	s = strings.TrimSuffix(s, ".rm")
	s = strings.Replace(s, "/@v/", " ", -1) // in web handlers
	s = strings.Replace(s, "@v", " v", -1)  // pasted from logs

	var mod coordinates.Module
	split := strings.Fields(s)
	if len(split) < 2 {
		return mod, errors.Errorf("malformed module line: %q", orig)
	}

	source := strings.TrimSuffix(split[0], "/")
	version := strings.TrimSuffix(split[1], "/go.mod")

	mod.Source = source
	mod.Version = version
	return mod, nil
}
