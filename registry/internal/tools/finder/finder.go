package finder

import (
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"gophers.dev/pkgs/loggy"
	"gophers.dev/pkgs/semantic"

	"oss.indeed.com/go/modprox/pkg/clients/zips"
)

type Result struct {
	Text   string
	Latest Head
	Tags   []semantic.Tag
}

type Head struct {
	// Pseudo represents Go's custom version string for SHAs which are
	// not represented by a SemVer string.
	// e.g.
	Custom string
	Commit string
}

type Tag struct {
	SemVer string
	Commit string
}

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i Versions -s _mock.go

type Versions interface {
	// Request the list of semver tags set in the source git repository.
	Request(source string) (*Result, error)
}

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i Finder -s _mock.go

type Finder interface {
	// Find returns the special form module name for the latest commit,
	// as well as a list of tags that follow proper semver format understood
	// by the Go compiler.
	Find(string) (*Result, error)
}

type Options struct {
	Timeout     time.Duration
	Versions    map[string]Versions
	ProxyClient zips.ProxyClient
}

func New(opts Options) Finder {
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = 1 * time.Minute
	}

	client := &http.Client{
		Timeout: timeout,
	}

	versions := opts.Versions
	if versions == nil {
		versions = map[string]Versions{
			"github.com": Github("", client, opts.ProxyClient),
		}
	}

	return &finder{
		versions: versions,
		log:      loggy.New("finder"),
	}
}

type finder struct {
	versions map[string]Versions
	log      loggy.Logger
}

func (f *finder) Find(source string) (*Result, error) {
	resolver, err := f.forSource(source)
	if err != nil {
		return nil, err
	}
	return resolver.Request(source)
}

func parseDomain(source string) string {
	split := strings.Split(source, "/")
	return split[0]
}

func (f *finder) forSource(source string) (Versions, error) {
	domain := parseDomain(source)
	versions, exists := f.versions[domain]
	if !exists {
		return nil, errors.Errorf("no version resolver for domain %q", domain)
	}
	return versions, nil
}

func Compatible(source string) bool {
	// as more things are added, add them here
	return githubPkgRe.MatchString(source)
}
