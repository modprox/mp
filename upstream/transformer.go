package upstream

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shoenig/toolkit"

	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/libmodprox/repository"
)

type Resolver interface {
	Resolve(repository.ModInfo) (*Request, error)
}

type Transform interface {
	Modify(*Request) *Request
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
	domains    map[string]bool // only implement redirect metadata
	httpClient *http.Client
	log        loggy.Logger
}

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

	return &GoGetTransform{
		domains: match,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		log: loggy.New("go-get-transform"),
	}
}

func (t *GoGetTransform) Modify(r *Request) *Request {
	if !t.domains[r.Domain] {
		t.log.Tracef("domain %s is not set for go-get redirects", r.Domain)
		return r
	}
	t.log.Infof("doing go-get redirect lookup for domain %s", r.Domain)

	meta, err := t.doGoGetRequest(r)
	if err != nil {
		t.log.Warnf("unable to lookup go get redirect to %s (assuming none): %v", meta, err)
		return r
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

	return modified
}

type goGetMeta struct {
	transport string
	domain    string
	path      string
}

func (t *GoGetTransform) doGoGetRequest(r *Request) (goGetMeta, error) {
	var meta goGetMeta
	uri := fmt.Sprintf("%s://%s/%s?go-get=1", r.Transport, r.Domain, strings.Join(r.Namespace, "/"))
	request, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return meta, err
	}

	response, err := t.httpClient.Do(request)
	if err != nil {
		return meta, err
	}
	defer toolkit.Drain(response.Body)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return meta, err
	}

	// fmt.Println("lll body:", string(body))

	return parseGoGetMetadata(string(body))
}

var (
	sourceRe = regexp.MustCompile(`(http[s]?)://([\w-.]+)/([\w-./]+)`)
)

// gives us transport, domain, path
func parseGoGetMetadata(content string) (goGetMeta, error) {
	var meta goGetMeta
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, `name="go-source"`) {
			groups := sourceRe.FindStringSubmatch(line)
			fmt.Println("iii groups:", groups)
			if len(groups) != 4 {
				return meta, errors.Errorf("malformed go-source meta tag: %q", line)
			}
			return goGetMeta{
				transport: groups[1],
				domain:    groups[2],
				path:      groups[3],
			}, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return meta, err
	}
	return meta, errors.New("no go-source meta tag in response")
}

type DomainPathTransform struct {
	pathFmt string
}

// e.g. https://github.com/shoenig/petrify/archive/v4.0.1.zip
// e.g. https://gitlab.com/cryptsetup/cryptsetup/-/archive/v2.0.1/cryptsetup-v2.0.1.zip

func (t *DomainPathTransform) Modify(r *Request) *Request {
	version := addressableVersion(r.Version) // this seems a little conflated
	newPath := formatPath(t.pathFmt, version, r.Namespace)
	return &Request{
		Transport: r.Transport,
		Domain:    r.Domain,
		Namespace: r.Namespace,
		Version:   r.Version,
		Path:      newPath,
	}
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

var DefaultPathTransforms = map[string]Transform{
	"github.com": NewDomainPathTransform("ELEM1/ELEM2/archive/VERSION.zip"),
	"gitlab.com": NewDomainPathTransform("ELEM1/ELEM2/-/archive/VERSION/ELEM2-VERSION.zip"),
	"":           NewDomainPathTransform(""), // unknown
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

func (t *SetPathTransform) Modify(r *Request) *Request {
	domainPathTransform := t.domainPathTransforms[r.Domain]
	modified := domainPathTransform.Modify(r)
	t.log.Tracef("original: %s", r)
	t.log.Tracef("modified: %s", modified)
	return modified
}
