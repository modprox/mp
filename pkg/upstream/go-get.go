package upstream

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/shoenig/httplus/responses"

	"oss.indeed.com/go/modprox/pkg/loggy"
)

var maxLoggedBody = 500

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
	defer responses.Drain(response)

	bs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return meta, err
	}

	code := response.StatusCode
	body := string(bs)

	if code >= 400 {
		t.log.Errorf("failed to do go-get redirect, received code %d from %s", code, uri)
		if len(body) <= maxLoggedBody {
			t.log.Errorf("response body: %s", body)
		} else {
			t.log.Errorf("response body: %s...", body[:maxLoggedBody])
		}
		return meta, errors.Errorf("bad response code (%d) from %s", code, uri)
	}

	return parseGoGetMetadata(body)
}

var (
	sourceRe = regexp.MustCompile(`(http[s]?)://([\w-.]+)/([\w-./]+)`)
	log      = loggy.New("go-get")
)

// gives us transport, domain, path
func parseGoGetMetadata(content string) (goGetMeta, error) {
	if ggm, exists, err := tryParseGoMetaTag("go-source", content); err != nil {
		return ggm, err
	} else if exists {
		log.Infof("found go-source tag: %#v", ggm)
		return ggm, nil
	}

	if ggm, exists, err := tryParseGoMetaTag("go-import", content); err != nil {
		return ggm, err
	} else if exists {
		log.Infof("found go-import tag %#v", ggm)
		return ggm, nil
	}

	return goGetMeta{}, errors.New("neither go-source or go-import meta tag found")
}

// ghetto hack where we look for go-source first, which is usually
// the true github.com source
//
// only when this does not work do we use the go-import line, which
// may redirect to a vcs protocol.
func tryParseGoMetaTag(tag, content string) (goGetMeta, bool, error) {
	content = formatContent(content) // pre-process html

	var meta goGetMeta
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		metaTag := fmt.Sprintf("name=%q", tag)
		if strings.Contains(line, metaTag) {
			groups := sourceRe.FindStringSubmatch(line)
			if len(groups) != 4 {
				return meta, false, errors.Errorf("malformed meta tag: %q", line)
			}
			return goGetMeta{
				transport: groups[1],
				domain:    groups[2],
				path:      cleanupPath(groups[3]),
			}, true, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return meta, false, err
	}
	return meta, false, nil
}

func cleanupPath(p string) string {
	a := strings.TrimSuffix(p, "/")
	b := strings.TrimSuffix(a, ".git")
	return b
}

// need to rewrite newlines not preceded by closing angle bracket to be spaces
func formatContent(content string) string {
	content = strings.Replace(content, "\n", " ", -1)
	content = strings.Replace(content, ">", ">\n", -1)
	return content
}
