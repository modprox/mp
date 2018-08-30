package upstream

import (
	"errors"
	"fmt"
	"strings"

	"github.com/modprox/libmodprox/loggy"
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

type Transform interface {
	Modify(*Request) *Request
}

type Namespace []string

type Request struct {
	Transport string
	Domain    string
	Namespace Namespace
	Version   string
	Path      string
}

func (r *Request) String() string {
	return fmt.Sprintf(
		"[%q %q %v %q %q]",
		r.Transport,
		r.Domain,
		r.Namespace,
		r.Version,
		r.Path,
	)
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

func (r *Request) URI() string {
	return fmt.Sprintf("%s://%s/%s", r.Transport, r.Domain, r.Path)
}

func NewRequest(mod repository.ModInfo) (*Request, error) {
	domain, namespace, err := splitSource(mod.Source)
	if err != nil {
		return nil, err
	}
	return &Request{
		Transport: "https",
		Domain:    domain,
		Namespace: namespace,
		Version:   mod.Version,
	}, nil
}

func splitSource(s string) (string, Namespace, error) {
	split := strings.Split(s, "/")

	if len(split) <= 1 {
		return "", nil, errors.New("source does not contain a path")
	}

	domain := split[0]
	namespace := Namespace(split[1:])
	return domain, namespace, nil
}

type RedirectTransform struct {
	original     string
	substitution string
	log          loggy.Logger
}

func NewRedirectTransform(original, substitution string) Transform {
	return &RedirectTransform{
		original:     original,
		substitution: substitution,
		log:          loggy.New("redirect-transform"),
	}
}

func (t *RedirectTransform) Modify(r *Request) *Request {
	newDomain := r.Domain
	if newDomain == t.original {
		newDomain = t.substitution
	}

	modified := &Request{
		Transport: r.Transport,
		Domain:    newDomain,
		Namespace: r.Namespace,
		Version:   r.Version,
	}

	t.log.Tracef("original: %s", r)
	t.log.Tracef("modified: %s", modified)
	return modified
}

type DomainPathTransform struct {
	pathFmt string
}

// e.g. https://github.com/shoenig/petrify/archive/v4.0.1.zip
// e.g. https://gitlab.com/cryptsetup/cryptsetup/-/archive/v2.0.1/cryptsetup-v2.0.1.zip

func (t *DomainPathTransform) Modify(r *Request) *Request {
	newPath := formatPath(t.pathFmt, r.Version, r.Namespace)
	return &Request{
		Transport: r.Transport,
		Domain:    r.Domain,
		Namespace: r.Namespace,
		Version:   r.Version,
		Path:      newPath,
	}
}

// e.g. ELEM1/ELEM2/archive/VERSION.zip => shoenig/petrify/archive/v4.0.1.zip
// e.g. ELEM1/ELEM2/-/archive/VERSION/ELEM2-VERSION.zip => crypo/cryptsetup/-/archive/v2.0.1/cryptsetup-v2.0.1.zip

func formatPath(pathFmt, version string, namespace Namespace) string {
	var path = pathFmt
	for i, elem := range namespace {
		elemIdx := fmt.Sprintf("ELEM%d", i+1)
		path = strings.Replace(path, elemIdx, elem, -1)
	}
	path = strings.Replace(path, "VERSION", version, -1)
	return path
}

func NewDomainPathTransform(pathFmt string) Transform {
	return &DomainPathTransform{
		pathFmt: pathFmt,
	}
}

var DefaultPathTransforms = map[string]Transform{
	"github.com": NewDomainPathTransform("ELEM1/ELEM2/archive/VERSION.zip"),
	"gitlab.com": NewDomainPathTransform("ELEM1/ELEM2/-/archive/VERSION/ELEM2-VERSION.zip"),
	"":           NewDomainPathTransform("ELEM1/ELEM2/VERSION.zip"), // arbitrary
}

type SetPathTransform struct {
	domainPathTransforms map[string]Transform
	log                  loggy.Logger
}

func NewSetPathTransform(customDomainPathTransforms map[string]Transform) Transform {
	combined := combinedDomainPathTransforms(customDomainPathTransforms)
	return &SetPathTransform{
		domainPathTransforms: combined,
		log:                  loggy.New("set-path-transform"),
	}
}

func combinedDomainPathTransforms(
	customDomainPathTransforms map[string]Transform,
) map[string]Transform {
	m := make(map[string]Transform, len(DefaultPathTransforms)+len(customDomainPathTransforms))
	for domain, transform := range DefaultPathTransforms {
		m[domain] = transform
	}
	for domain, transform := range customDomainPathTransforms {
		m[domain] = transform
	}
	return m
}

func (s *SetPathTransform) Modify(r *Request) *Request {
	domainPathTransform := s.domainPathTransforms[r.Domain]
	modified := domainPathTransform.Modify(r)
	s.log.Tracef("original: %s", r)
	s.log.Tracef("modified: %s", modified)
	return modified
}
