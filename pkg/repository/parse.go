package repository

import (
	"strings"

	"github.com/modprox/mp/pkg/coordinates"

	"github.com/pkg/errors"
)

var (
// examples
//  mod file style
//   github.com/foo/bar v2.0.0
//   github.com/tdewolff/parse v2.3.3+incompatible // indirect
//   golang.org/x/tools v0.0.0-20180111040409-fbec762f837d
//   gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405
//  proxy request style
//   /github.com/cpuguy83/go-md2man/@v/v1.0.6.info
)

// Parse will parse s as a module in string form.
func Parse(s string) (coordinates.Module, error) {
	orig := s
	s = strings.Trim(s, "/")
	s = strings.TrimSuffix(s, ".info")
	s = strings.TrimSuffix(s, ".zip")
	s = strings.TrimSuffix(s, ".mod")
	s = strings.Replace(s, "/@v/", " ", -1)

	var mod coordinates.Module
	split := strings.Fields(s)
	if len(split) < 2 {
		return mod, errors.Errorf("malformed module line: %q", orig)
	}

	source := strings.TrimSuffix(split[0], "/")
	version := split[1]

	mod.Source = source
	mod.Version = version
	return mod, nil
}
