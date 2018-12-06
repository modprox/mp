package web

import (
	"strings"

	"github.com/modprox/mp/pkg/coordinates"
	"github.com/modprox/mp/pkg/repository"
)

// GET baseURL/module/@v/list fetches a list of all known versions, one per line.

func moduleFromPath(p string) (string, error) {
	p = demangle(p)
	mod, err := repository.Parse(p)
	return mod.Source, err
}

func modInfoFromPath(p string) (coordinates.Module, error) {
	p = demangle(p)
	return repository.Parse(p)
}

// from the Go documentation: https://tip.golang.org/cmd/go/#hdr-Module_proxy_protocol
//
// To avoid problems when serving from case-sensitive file systems, the <module> and <version>
// elements are case-encoded, replacing every uppercase letter with an exclamation mark followed
// by the corresponding lower-case letter: github.com/Azure encodes as github.com/!azure.
//
// modprox currently store modules under their correct names, so we must rewrite the go
// commands download requests from the mangled name to the correct name.
func demangle(s string) string {
	var correct strings.Builder

	// copy s into correct, using lookahead to rewrite letters
	var i int
	for i = 0; i < len(s)-1; i++ {
		c := s[i]
		if c == '!' {
			next := s[i+1]
			if next >= 'a' && next <= 'z' {
				correct.WriteByte(next - ('a' - 'A'))
				i++
				continue
			}
		}
		correct.WriteByte(c)
	}

	// if the last 2 letters were not an encoding, copy the last letter
	if i == len(s)-1 {
		correct.WriteByte(s[i])
	}

	return correct.String()
}
