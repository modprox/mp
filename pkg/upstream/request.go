package upstream

import (
	"fmt"
	"strings"
)

// Namespace is the path elements leading up to the package name of a module.
type Namespace []string

// A Request for a git-archive of a module.
//
// A Request should only be used in the context of acquiring a module archive
// from an upstream source (e.g. github, gitlab, etc.) as opposed to acquiring
// a module from a proxy.
type Request struct {
	Transport     string
	Domain        string
	Namespace     Namespace
	Version       string
	Path          string
	GoGetRedirect bool
	Headers       map[string]string
}

func (r *Request) String() string {
	return fmt.Sprintf(
		"[%q %q %v %q %q %t]",
		r.Transport,
		r.Domain,
		r.Namespace,
		r.Version,
		r.Path,
		r.GoGetRedirect,
	)
}

// An explicit implementation of equality between two Request objects.
func (r *Request) Equals(o *Request) bool {
	if r.Transport != o.Transport {
		return false
	}

	if r.Domain != o.Domain {
		return false
	}

	if len(r.Namespace) != len(o.Namespace) {
		return false
	}

	for i := 0; i < len(r.Namespace); i++ {
		if r.Namespace[i] != o.Namespace[i] {
			return false
		}
	}

	if r.Version != o.Version {
		return false
	}

	if r.Path != o.Path {
		return false
	}

	if r.GoGetRedirect != o.GoGetRedirect {
		return false
	}

	if len(r.Headers) != len(o.Headers) {
		return false
	}

	for key := range r.Headers {
		if r.Headers[key] != o.Headers[key] {
			return false
		}
	}

	return true
}

// The URI is only valid AFTER a Request has passed through
// all of the Transform functors.
//
// The URI should represent the way to get some git-archive, which must
// later be transformed into a proper module archive before the Go tooling
// will be able to work with it.
func (r *Request) URI() string {
	rPath := strings.TrimPrefix(r.Path, "/")
	return fmt.Sprintf("%s://%s/%s", r.Transport, r.Domain, rPath)
}
