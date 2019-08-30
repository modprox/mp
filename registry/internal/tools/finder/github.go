package finder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gophers.dev/pkgs/ignore"
	"gophers.dev/pkgs/loggy"
)

func Github(baseURL string, client *http.Client) Versions {
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}
	return &github{
		baseURL: baseURL,
		client:  client,
		log:     loggy.New("github-versions"),
	}
}

type github struct {
	baseURL string
	client  *http.Client
	log     loggy.Logger
}

func (g *github) Request(source string) (*Result, error) {
	namespace, project, err := g.parseSource(source)
	if err != nil {
		return nil, err
	}

	tagsURI := g.tagsURI(namespace, project)

	g.log.Tracef("requesting tags from URI: %s", tagsURI)

	tags, err := g.requestTags(tagsURI)
	if err != nil {
		return nil, err
	}

	headURI := g.headURI(namespace, project)

	g.log.Tracef("requesting latest commit from URI: %s", headURI)

	head, err := g.requestHead(headURI, tags)
	if err != nil {
		return nil, err
	}

	return &Result{
		Latest: head,
		Tags:   tags,
	}, nil
}

func (g *github) requestTags(uri string) ([]Tag, error) {
	response, err := g.client.Get(uri)
	if err != nil {
		return nil, err
	}
	defer ignore.Drain(response.Body)

	return g.decodeTags(response.Body)
}

func (g *github) requestHead(uri string, tags []Tag) (Head, error) {
	response, err := g.client.Get(uri)
	if err != nil {
		return Head{}, err
	}
	defer ignore.Drain(response.Body)

	return g.decodeHead(response.Body, tags)
}

func (g *github) decodeHead(r io.Reader, tags []Tag) (Head, error) {
	var gCommit githubCommit
	if err := json.NewDecoder(r).Decode(&gCommit); err != nil {
		return Head{}, err
	}

	custom, err := gCommit.Pseudo(tags)
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

func (gc githubCommit) Pseudo(tags []Tag) (string, error) {
	ts, err := time.Parse(time.RFC3339, gc.Commit.Author.Date)
	if err != nil {
		return "", err
	}

	date := ts.Format("200601020300")
	shortSHA := gc.SHA[0:12] // what Go does

	if len(tags) == 0 {
		return fmt.Sprintf("v0.0.0-%s-%s", date, shortSHA), nil
	}

	lastVersion := tags[0].SemVer
	if strings.HasSuffix(lastVersion, "-pre") {
		return fmt.Sprintf("%s.0.%s-%s", lastVersion, date, shortSHA), nil
	}

	semver := parseSemVer(lastVersion)
	if semver == nil {
		// illegal semver
		return fmt.Sprintf("v0.0.0-%s-%s", date, shortSHA), nil
	}
	return fmt.Sprintf("v%d.%d.%d-0.%s-%s", semver.Major, semver.Minor, semver.Patch+1, date, shortSHA), nil
}

var semVerRe = regexp.MustCompile(`^v(\d+)(?:\.(\d+)(?:\.(\d+))?)?$`)

func parseSemVer(semver string) *SemVer {
	matches := semVerRe.FindStringSubmatch(semver)
	if matches == nil {
		return &SemVer{0, 0, 0}
	}
	if matches[2] == "" && matches[3] == "" {
		return &SemVer{unsafeIToA(matches[1]), 0, 0}
	}
	if matches[3] == "" {
		return &SemVer{unsafeIToA(matches[1]), unsafeIToA(matches[2]), 0}
	}
	return &SemVer{unsafeIToA(matches[1]), unsafeIToA(matches[2]), unsafeIToA(matches[3])}
}

// only call this if you are sure that s is convertible to an int
func unsafeIToA(s string) int {
	res, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Errorf("failed to convert %s to an int ; this should never happen", s))
	}
	return res
}

func (g *github) decodeTags(r io.Reader) ([]Tag, error) {
	var gTags []githubTag
	if err := json.NewDecoder(r).Decode(&gTags); err != nil {
		return nil, err
	}
	var tags []Tag
	for _, gTag := range gTags {
		tags = append(tags, Tag{
			SemVer: gTag.Name,
			Commit: gTag.Commit.SHA,
		})
	}
	return tags, nil
}

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

func (g *github) tagsURI(namespace, project string) string {
	return fmt.Sprintf(
		"%s/repos/%s/%s/tags",
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

// because we need to parse githubs response
type githubTag struct {
	Name   string `json:"name"`
	Commit struct {
		SHA string `json:"sha"`
	} `json:"commit"`
}
