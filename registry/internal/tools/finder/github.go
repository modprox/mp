package finder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/shoenig/httplus/responses"

	"go.gophers.dev/pkgs/loggy"
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
	tagsURI, err := g.tagsURI(source)
	if err != nil {
		return nil, err
	}

	g.log.Tracef("requesting tags from URI: %s", tagsURI)

	tags, err := g.requestTags(tagsURI)
	if err != nil {
		return nil, err
	}

	headURI, err := g.headURI(source)
	if err != nil {
		return nil, err
	}

	g.log.Tracef("requesting latest commit from URI: %s", headURI)

	head, err := g.requestHead(headURI)
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
	defer responses.Drain(response)

	return g.decodeTags(response.Body)
}

func (g *github) requestHead(uri string) (Head, error) {
	response, err := g.client.Get(uri)
	if err != nil {
		return Head{}, err
	}
	defer responses.Drain(response)

	return g.decodeHead(response.Body)
}

func (g *github) decodeHead(r io.Reader) (Head, error) {
	var gCommit githubCommit
	if err := json.NewDecoder(r).Decode(&gCommit); err != nil {
		return Head{}, err
	}

	custom, err := gCommit.Pseudo()
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

func (gc githubCommit) Pseudo() (string, error) {
	ts, err := time.Parse(time.RFC3339, gc.Commit.Author.Date)
	if err != nil {
		return "", err
	}

	date := ts.Format("200601020300")

	pseudo := fmt.Sprintf(
		"v0.0.0-%s-%s+incompatible",
		date,
		gc.SHA[0:12], // what Go does
	)

	return pseudo, nil
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

func (g *github) headURI(source string) (string, error) {
	namespace, project, err := g.parseSource(source)
	if err != nil {
		return "", err
	}

	headURI := fmt.Sprintf(
		"%s/repos/%s/%s/commits/HEAD",
		g.baseURL,
		namespace,
		project,
	)

	return headURI, nil
}

func (g *github) tagsURI(source string) (string, error) {
	namespace, project, err := g.parseSource(source)
	if err != nil {
		return "", err
	}

	apiURI := fmt.Sprintf(
		"%s/repos/%s/%s/tags",
		g.baseURL,
		namespace,
		project,
	)

	return apiURI, nil
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
