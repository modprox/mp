package upstream

import (
	"errors"
	"fmt"
	"strings"

	"github.com/modprox/libmodprox/repository"
)

// what we really need is a thing which
// transforms a module into a URI usable for an http
// request - by applying each of the types of transforms:
//
// - domain alias
// - URL path creation based on domain
// - authentication / authorization configuration

type Resolver interface {
	Resolve(repository.ModInfo) (*Request, error)
}

type resolver struct {
	transforms []Transform
}

func NewResolver(transforms ...Transform) Resolver {
	return &resolver{
		transforms: transforms,
	}
}

func (r *resolver) Resolve(mod repository.ModInfo) (*Request, error) {
	request, err := NewRequest(mod)
	if err != nil {
		return nil, err
	}

	for _, transform := range r.transforms {
		request = transform.Modify(request)
	}

	return request, nil
}

type Transform interface {
	Modify(*Request) *Request
}

type Request struct {
	Transport string
	Domain    string
	Path      string
	Version   string
}

func (r *Request) URI() string {
	return fmt.Sprintf("%s://%s/%s", r.Transport, r.Domain, r.Path)
}

func NewRequest(mod repository.ModInfo) (*Request, error) {
	domain, path, err := splitSource(mod.Source)
	if err != nil {
		return nil, err
	}
	return &Request{
		Transport: "https",
		Domain:    domain,
		Path:      path,
		Version:   mod.Version,
	}, nil
}

func splitSource(s string) (string, string, error) {
	split := strings.SplitN(s, "/", 2)
	if len(split) != 2 {
		return "", "", errors.New("source does not contain a path")
	}
	return split[0], split[1], nil
}

type RedirectTransform struct {
	original     string
	substitution string
}

func NewRedirectTransform(original, substitution string) Transform {
	return &RedirectTransform{
		original:     original,
		substitution: substitution,
	}
}

func (t *RedirectTransform) Modify(r *Request) *Request {
	newDomain := r.Domain
	if newDomain == t.original {
		newDomain = t.substitution
	}

	return &Request{
		Transport: r.Transport,
		Domain:    newDomain,
		Path:      r.Path,
		Version:   r.Version,
	}
}
