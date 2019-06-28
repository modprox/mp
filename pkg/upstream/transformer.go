package upstream

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"oss.indeed.com/go/modprox/pkg/coordinates"
	"oss.indeed.com/go/modprox/pkg/loggy"
)

//go:generate go run github.com/gojuno/minimock/cmd/minimock -g -i Resolver -s _mock.go

// A Resolver is able to turn the globally unique identifier of
// a Go module (which includes a Source and a Version) and applies
// a set of Transform operations until a Request is created
// that can later be used to fetch the module from some source,
// which is typically a VCS host (e.g. github).
type Resolver interface {
	// Resolve applies any underlying Transform operations
	// and returns the resulting Request, or an error if
	// one of the Transform operations does not work.
	Resolve(coordinates.Module) (*Request, error)

	// UseProxy indicates whether Module can be downloaded from a generic open
	// global proxy (e.g. proxy.golang.org) instead of the original upstream
	// source (e.g. github / gitlab). This is going to be true for any Module
	// which does not get matched by any configured domain type Transform. The
	// rational is that any domain type Transform is a flag that the module is
	// not going to be present an the open source context, and the original
	// upstream must be used since it is likely a private repository.
	//
	// The transforms that prohibit proxy use are:
	// - StaticRedirectTransform
	// - DomainTransportTransform
	// - DomainHeaderTransform
	UseProxy(coordinates.Module) (bool, error)
}

// A Transform is one operation that is applied to a Request,
// which creates a new Request with zero or more parameters
// of the input Request having been modified. A Transform can
// be used to handle things like static domain name redirection,
// indirect domain name redirect (i.e. accommodate go-get meta URIs),
// domain-based path rewriting, etc.
//
// As time goes on, more and more Transform implementations will be
// added, to support additional use cases for enterprise environments
// which tend to have special needs.
type Transform interface {
	Modify(*Request) (*Request, error)
}

type resolver struct {
	transforms []Transform
}

// NewResolver creates a Resolver which will apply the given set
// of Transform operations in the order in which they appear.
func NewResolver(transforms ...Transform) Resolver {
	return &resolver{
		transforms: transforms,
	}
}

func (r *resolver) UseProxy(mod coordinates.Module) (bool, error) {
	original, err := NewRequest(mod)
	if err != nil {
		return false, err
	}

	// go through each transform and decide if it applies to this module
	for _, transform := range r.transforms {

		switch transform.(type) {

		// select on the types that trigger an upstream request to be necessary
		case *StaticRedirectTransform,
			*DomainTransportTransform,
			*DomainHeaderTransform:

			// Apply the transform, if the request is modified, that means
			// the transform is applies to this module and we cannot use the
			// global proxy to make the request.
			changed, err := transform.Modify(original)
			if err != nil {
				return false, err
			}

			// Compare the original request with the modified request. If they
			// do not match, we cannot use a global proxy to make request the
			// archive for this module.
			if !original.Equals(changed) {
				return false, nil
			}
		}
	}

	return true, nil
}

func (r *resolver) Resolve(mod coordinates.Module) (*Request, error) {
	request, err := NewRequest(mod)
	if err != nil {
		return nil, err
	}

	for _, transform := range r.transforms {
		request, err = transform.Modify(request)
		if err != nil {
			return nil, err
		}
	}

	return request, nil
}

// NewRequest creates a default Request from the given module. This
// initial Request is likely useless, as it only becomes useful after
// a set of Transform operations are applied to it, which then compute
// correct URI for the module it represents.
func NewRequest(mod coordinates.Module) (*Request, error) {
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
	s = strings.TrimSpace(s)
	if s == "" {
		return "", nil, errors.New("source is empty string")
	}

	split := strings.Split(s, "/")

	if len(split) == 1 {
		// e.g. go.opencensus.io is a whole domain used
		// to represent one package using go-get meta
		return s, nil, nil
	}

	// we have a domain and a path

	domain := split[0]
	namespace := Namespace(split[1:])
	return domain, namespace, nil
}

// A StaticRedirectTransform is used to directly convert one domain
// to another. For example, if your organization internally keeps packages
// organized like
//   ${GOPATH}/company/...
// but the internal VCS is only addressable in a way like
//   code.internal.company.net/...
// then the StaticRedirectTransform can be used to automatically acquire
// modules prefixed with name "company/" from the internal VCS of the
// different domain name.
type StaticRedirectTransform struct {
	original     string
	substitution string
	log          loggy.Logger
}

// NewStaticRedirectTransform creates a Transform which will convert
// domains of the original name to become the substitution name.
//
// Currently only exact matches on the domain are supported.
func NewStaticRedirectTransform(original, substitution string) Transform {
	return &StaticRedirectTransform{
		original:     original,
		substitution: substitution,
		log:          loggy.New("redirect-transform"),
	}
}

func (t *StaticRedirectTransform) Modify(r *Request) (*Request, error) {
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
	return modified, nil
}

// The GoGetTransform triggers an http request to the domain
// to simply do a "?go-get=1" lookup for the real domain of where
// the module is being hosted.
//
// Additional domains can be specified via configuration.
// The known go-get redirectors in the wild include:
// - golang.org
// - google.golang.org
// - cloud.google.com
// - gopkg.in
// - contrib.go.opencensus.io
// - go.uber.org
type GoGetTransform struct {
	autoRedirect bool
	domains      map[string]bool // only implement redirect metadata
	httpClient   *http.Client
	log          loggy.Logger
}

// NewAutomaticGoGetTransform creates a GoGetTransform where any module URI
// will be redirected to wherever the go-get meta HTML tag in the domain
// indicates.
func NewAutomaticGoGetTransform() Transform {
	return &GoGetTransform{
		autoRedirect: true,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		log: loggy.New("auto-go-get-transform"),
	}
}

// NewGoGetTransform creates a GoGetTransform where any module URIs
// found in the given list of domains will be first redirected to wherever
// the go-get meta HTML tag in the domain indicates.
//
// Read more about this functionality here:
//   https://golang.org/cmd/go/#hdr-Remote_import_paths
func NewGoGetTransform(domains []string) Transform {
	match := make(map[string]bool)
	for _, domain := range domains {
		match[domain] = true
	}

	match["golang.org"] = true
	match["cloud.google.com"] = true
	match["google.golang.org"] = true
	match["gopkg.in"] = true
	match["contrib.go.opencensus.io"] = true
	match["go.opencensus.io"] = true
	match["go.uber.org"] = true
	match["git.apache.org"] = true
	match["k8s.io"] = true
	match["sigs.k8s.io"] = true

	return &GoGetTransform{
		domains: match,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		log: loggy.New("go-get-transform"),
	}
}

func (t *GoGetTransform) Modify(r *Request) (*Request, error) {
	if !t.autoRedirect && !t.domains[r.Domain] {
		t.log.Tracef("domain %s is not set for go-get redirects", r.Domain)
		return r, nil
	}

	t.log.Infof("doing go-get redirect lookup for domain %s", r.Domain)

	meta, err := t.doGoGetRequest(r)
	if !t.autoRedirect && err != nil {
		t.log.Errorf("unable to do go-get redirect for domain %s: %v", r.Domain, err)
		return nil, err
	} else if err != nil {
		t.log.Warnf("unable to do go-get redirect for domain: %s. leaving request unmodified: %v", r.Domain, err)
		return r, nil
	}

	t.log.Infof("redirect to: %s", meta)
	modified := &Request{
		Transport: meta.transport,
		Domain:    meta.domain,
		Namespace: strings.Split(meta.path, "/"),
		Version:   r.Version,
		// Path: set by the domain rewriter
	}

	t.log.Tracef("original: %s", r)
	t.log.Tracef("modified: %s", modified)

	return modified, nil
}

// A DomainPathTransform is used to generate or rewrite the URL path
// of the module archive that is to be fetched per the domain of desired
// module of the Request. Default path rewriting rules are provided for
// repositories ultimately hosted in github or gitlab. Additional path
// transformations should be defined for internally hosed VCSs.
//
// e.g. github:
//   https://github.com/ELEM1/ELEM2/archive/VERSION.zip
// e.g. gitlab:
//   https://gitlab.com/ELEM1/ELEM2/-/archive/VERSOIN/ELEM2-v2.0.1.zip
type DomainPathTransform struct {
	pathFmt string
}

func (t *DomainPathTransform) Modify(r *Request) (*Request, error) {
	version := addressableVersion(r.Version) // this seems a little conflated
	newPath := formatPath(t.pathFmt, version, r.Namespace)
	return &Request{
		Transport: r.Transport,
		Domain:    r.Domain,
		Namespace: r.Namespace,
		Version:   r.Version,
		Path:      newPath,
	}, nil
}

// e.g. v2.0.0 => v2.0.0
// e.g. v0.0.0-20180111040409-fbec762f837d => fbec762f837d
// e.g. v2.3.3+incompatible => v2.3.3
func addressableVersion(version string) string {
	// dashes indicate <version>-<timestamp>-<hash> format,
	// where the hash is what is addressable in vcs
	splitOnDash := strings.Split(version, "-")
	if len(splitOnDash) == 3 {
		return splitOnDash[2] // return the hash if it exists
	}

	// plus indicates <version>+<comment> where the version
	// is what is addressable in vcs
	splitOnPlus := strings.Split(version, "+")
	if len(splitOnPlus) > 1 {
		return splitOnPlus[0]
	}

	// the version is just the version
	return version
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

// DefaultPathTransforms provides a set of default Transform types which
// create the Request.Path for a known set of VCSs systems in the open source
// world (i.e. github and gitlab).
// Additional Transforms should be specified via NewSetPathTransform, which
// accepts a map of domain to Transform, for internally hosed code.
var DefaultPathTransforms = map[string]Transform{
	"github.com": NewDomainPathTransform("ELEM1/ELEM2/archive/VERSION.zip"),
	"gitlab.com": NewDomainPathTransform("ELEM1/ELEM2/-/archive/VERSION/ELEM2-VERSION.zip"),
	"":           NewDomainPathTransform(""), // unknown
}

// A SetPathTransform is a collection of transforms which set the Path of
// a Request given a domain. Think of it as a map from a domain to a
// DomainPathTransform, which can be used in the general case rather than
// specifying an explicit list of DomainPathTransform.
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

func (t *SetPathTransform) Modify(r *Request) (*Request, error) {
	domainPathTransform, exists := t.domainPathTransforms[r.Domain]
	if !exists {
		return nil, errors.Errorf("no path transformation exists for domain %s", r.Domain)
	}
	modified, err := domainPathTransform.Modify(r)
	t.log.Tracef("original: %s", r)
	t.log.Tracef("modified: %s", modified)
	return modified, err
}

func NewDomainHeaderTransform(domain string, headers map[string]string) Transform {
	return &DomainHeaderTransform{
		domain:  domain,
		headers: headers,
		log:     loggy.New("domain-header-transform"),
	}
}

// A DomainHeaderTransform is used to set the header for a request.
// Typically one of these will be used to set the authentication key
// for https requests to an internal VCS system.
type DomainHeaderTransform struct {
	domain  string
	headers map[string]string
	log     loggy.Logger
}

func (t *DomainHeaderTransform) Modify(r *Request) (*Request, error) {
	if r.Domain != t.domain {
		return r, nil
	}

	newHeaders := make(map[string]string, len(r.Headers))
	for k, v := range r.Headers {
		newHeaders[k] = v
	}

	for key, value := range t.headers {
		t.log.Tracef("setting a value for request header %q", key)
		newHeaders[key] = value
	}

	return &Request{
		Transport: r.Transport,
		Domain:    r.Domain,
		Namespace: r.Namespace,
		Version:   r.Version,
		Path:      r.Path,
		Headers:   newHeaders,
	}, nil
}

func NewDomainTransportTransform(domain, transport string) Transform {
	return &DomainTransportTransform{
		domain:    domain,
		transport: transport,
		log:       loggy.New("domain-transport-transform"),
	}
}

type DomainTransportTransform struct {
	domain    string
	transport string // e.g. https/http
	log       loggy.Logger
}

func (t *DomainTransportTransform) Modify(r *Request) (*Request, error) {
	if r.Domain != t.domain {
		return r, nil
	}

	newTransport := t.transport
	t.log.Tracef("setting transport of request to %q", newTransport)

	return &Request{
		Transport: newTransport,
		Domain:    r.Domain,
		Namespace: r.Namespace,
		Version:   r.Version,
		Path:      r.Path,
		Headers:   r.Headers,
	}, nil
}
