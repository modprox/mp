package finder

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"

	"gophers.dev/pkgs/ignore"
	"gophers.dev/pkgs/loggy"
	"gophers.dev/pkgs/semantic"

	"oss.indeed.com/go/modprox/pkg/clients/zips"
)

func Github(baseURL string, client *http.Client, proxyClient zips.ProxyClient) Versions {
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}
	return &github{
		baseURL:     baseURL,
		client:      client,
		log:         loggy.New("github-versions"),
		proxyClient: proxyClient,
	}
}

type github struct {
	baseURL     string
	client      *http.Client
	log         loggy.Logger
	proxyClient zips.ProxyClient
}

func (g *github) Request(source string) (*Result, error) {
	namespace, project, err := g.parseSource(source)
	if err != nil {
		return nil, err
	}

	g.log.Tracef("requesting available versions from the official go proxy")

	tags, err := g.proxyClient.List(source)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query list of versions from proxy.golang.org")
	}

	g.log.Tracef("checking if %s is module-compatible", source)
	isModuleCompatible, err := g.isModuleCompatible(tags)
	if err != nil {
		return nil, err
	}

	headURI := g.headURI(namespace, project)

	g.log.Tracef("requesting latest commit from URI: %s", headURI)

	head, err := g.requestHead(headURI, tags, isModuleCompatible)
	if err != nil {
		return nil, err
	}

	return &Result{
		Latest: head,
		Tags:   tags,
	}, nil
}

// isModuleCompatible isn't 100% accurate. It will only return false if the latest semver
// major >= 2
// This is good enough for the usecase intended, which is to decide on adding "+incompatible"
// to version strings. "+incompatible" isn't needed for major versions < 2
func (g *github) isModuleCompatible(versions []semantic.Tag) (bool, error) {
	return len(versions) > 0 && strings.Contains(versions[len(versions)-1].Extension, "+incompatible"), nil
}

func (g *github) requestHead(uri string, tags []semantic.Tag, isModuleCompatible bool) (Head, error) {
	response, err := g.client.Get(uri)
	if err != nil {
		return Head{}, err
	}
	defer ignore.Drain(response.Body)

	return g.decodeHead(response.Body, tags, isModuleCompatible)
}

func (g *github) decodeHead(r io.Reader, tags []semantic.Tag, isModuleCompatible bool) (Head, error) {
	var gCommit githubCommit
	if err := json.NewDecoder(r).Decode(&gCommit); err != nil {
		return Head{}, err
	}

	custom, err := gCommit.Pseudo(tags, isModuleCompatible)
	if err != nil {
		return Head{}, err
	}

	return Head{
		Commit: gCommit.SHA,
		Custom: custom,
	}, nil
}

type githubCommit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Author struct {
			Date string `json:"date"`
		} `json:"author"`
	} `json:"commit"`
}

func (gc githubCommit) Pseudo(tags []semantic.Tag, isModuleCompatible bool) (string, error) {
	naked, semver, err := gc.nakedPseudo(tags)
	if err != nil {
		return "", err
	}

	if semver == nil || semver.Major < 2 || isModuleCompatible {
		return naked, nil
	}
	return naked + "+incompatible", nil
}

func (gc githubCommit) nakedPseudo(tags []semantic.Tag) (string, *semantic.Tag, error) {
	ts, err := time.Parse(time.RFC3339, gc.Commit.Author.Date)
	if err != nil {
		return "", nil, err
	}

	date := ts.Format("20060102150405")
	shortSHA := gc.SHA[0:12] // what Go does

	if len(tags) == 0 {
		return fmt.Sprintf("v0.0.0-%s-%s", date, shortSHA), nil, nil
	}

	// tags are guaranteed to be logically reverse-ordered by proxyClient
	semver := tags[0]

	if semver.Extension == "pre" {
		// TODO should this always be ".0" or do we increment the pre version?
		return fmt.Sprintf("%s.0.%s-%s", semver.String(), date, shortSHA), &semver, nil
	}

	return fmt.Sprintf("v%d.%d.%d-0.%s-%s", semver.Major, semver.Minor, semver.Patch+1, date, shortSHA), &semver, nil
}

// -rc is commonly used, but not in the spec
//var semVerRe = regexp.MustCompile(`^v(\d+)(?:\.(\d+)(?:\.(\d+(-pre|-rc)?))?)?$`)

// only github.com things are supported for now
var githubPkgRe = regexp.MustCompile(`(github\.com)/([[:alnum:]_-]+)/([[:alnum:]_-]+)`)

func (g *github) headURI(namespace, project string) string {
	return fmt.Sprintf(
		"%s/repos/%s/%s/commits/HEAD",
		g.baseURL,
		namespace,
		project,
	)
}

func (g *github) parseSource(source string) (string, string, error) {
	groups := githubPkgRe.FindStringSubmatch(source)
	if len(groups) != 4 {
		return "", "", errors.New("source does not conform to format")
	}

	if groups[1] != "github.com" {
		return "", "", errors.New("only github.com is currently supported")
	}

	return groups[2], groups[3], nil
}
