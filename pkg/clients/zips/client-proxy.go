package zips

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	clean "github.com/hashicorp/go-cleanhttp"

	"github.com/pkg/errors"

	"gophers.dev/pkgs/ignore"
	"gophers.dev/pkgs/loggy"
	"gophers.dev/pkgs/semantic"

	"oss.indeed.com/go/modprox/pkg/coordinates"
	"oss.indeed.com/go/modprox/pkg/repository"
)

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i ProxyClient -s _mock.go

// A ProxyClient is used for making requests to a Go Module Proxy
// which is expected to return archives already in the correct format.
type ProxyClient interface {
	// Get returns the contents of the repo specified by the coordinates
	Get(coordinates.Module) (repository.Blob, error)
	// List returns all available versions of the repo specified by the coordinates, in descending logical order
	List(source string) ([]semantic.Tag, error)
}

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i iHTTPClient
type iHTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type proxyClient struct {
	httpClient iHTTPClient
	baseURL    string
	protocol   string
	log        loggy.Logger
}

type ProxyClientOptions struct {
	Protocol string        // typically https
	BaseURL  string        // typically proxy.golang.org
	Timeout  time.Duration // about 1 minute is good
}

// NewProxyClient creates a ProxyClient with some options.
func NewProxyClient(opts ProxyClientOptions) ProxyClient {
	if opts.BaseURL == "" {
		panic("proxy client BaseURL must be provided")
	}

	if opts.Protocol == "" {
		panic("proxy client protocol must be set")
	}

	if opts.Timeout <= 0 {
		panic("proxy client Timeout must be positive")
	}

	httpClient := clean.DefaultPooledClient()
	httpClient.Timeout = opts.Timeout

	return &proxyClient{
		baseURL:    opts.BaseURL,
		httpClient: httpClient,
		protocol:   opts.Protocol,
		log:        loggy.New("proxy-client"),
	}
}

func (c *proxyClient) zipURIOf(module coordinates.Module) string {
	modZipPath := mangle(fmt.Sprintf(
		"/%s/@v/%s.zip",
		module.Source,
		module.Version,
	))

	s := url.URL{
		Scheme: c.protocol,
		Host:   c.baseURL,
		Path:   modZipPath,
	}

	return s.String()
}

func (c *proxyClient) listURIOf(source string) string {
	modListPath := mangle(fmt.Sprintf("/%s/@v/list", source))

	s := url.URL{
		Scheme: c.protocol,
		Host:   c.baseURL,
		Path:   modListPath,
	}

	return s.String()
}

// module name needs to be re-encoded per the custom way Go decided to
// encode module names (read https://go.indeed.com/GLYM8RA6E)
//
// The safe encoding is this:
// replace every uppercase letter with an exclamation mark
// followed by the letter's lowercase equivalent.
//
// For example,
// github.com/Azure/azure-sdk-for-go ->  github.com/!azure/azure-sdk-for-go.
// github.com/GoogleCloudPlatform/cloudsql-proxy -> github.com/!google!cloud!platform/cloudsql-proxy
// github.com/Sirupsen/logrus -> github.com/!sirupsen/logrus.
//
// The opposite of this is in proxy/internal/web/common.go
// We could extract these into a library in the future.

func mangle(source string) string {
	var builder strings.Builder

	for _, c := range source {
		if c >= 'A' && c <= 'Z' {
			builder.WriteString("!")
			builder.WriteString(string(c + ('a' - 'A'))) // remember CS101?
		} else {
			builder.WriteString(string(c))
		}
	}

	return builder.String()
}

func (c *proxyClient) Get(mod coordinates.Module) (repository.Blob, error) {
	// request looks like
	//
	// GET https://proxy.golang.org/oss.indeed.com/go/taggit/@v/v0.3.3.zip
	zipURI := c.zipURIOf(mod)
	c.log.Tracef("making zip proxy request to %s", zipURI)

	response, err := c.sendRequest(mod.String(), zipURI)
	if err != nil {
		return nil, err
	}
	defer ignore.Drain(response)

	return ioutil.ReadAll(response)
}

func (c *proxyClient) List(source string) ([]semantic.Tag, error) {
	// request looks like
	//
	// GET https://proxy.golang.org/oss.indeed.com/go/taggit/@v/list
	listURI := c.listURIOf(source)
	c.log.Tracef("making list proxy request to %s", listURI)

	response, err := c.sendRequest(source, listURI)
	if err != nil {
		return nil, err
	}
	defer ignore.Drain(response)

	result := []semantic.Tag{}
	scanner := bufio.NewScanner(response)
	for scanner.Scan() {
		version := scanner.Text()
		tag, success := semantic.Parse(version)
		if !success {
			return nil, errors.Errorf("failed to parse version %s", version)
		}
		result = append(result, tag)
	}
	sort.Sort(sort.Reverse(semantic.BySemver(result)))

	return result, nil
}

func (c *proxyClient) sendRequest(subject, uri string) (io.ReadCloser, error) {
	// create the request for the module, from the proxy
	request, err := c.newRequest(uri)
	if err != nil {
		return nil, err
	}

	// do the request for the module, from the proxy
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrapf(err, "proxy request failed for %s", subject)
	}

	// if we get a bad response code, try to read the body and log it
	// todo: can we make this generic? copied from http.go
	if response.StatusCode >= 400 {
		bs, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, errors.Wrapf(err, "could not read body of bad response (%d)", response.StatusCode)
		}
		body := string(bs)
		if len(body) <= maxLoggedBody {
			c.log.Errorf("bad response (%d) body: %s", response.StatusCode, body)
		} else {
			c.log.Errorf("bad response(%d) trunc body: %s...", response.StatusCode, body[:maxLoggedBody])
		}
		return nil, errors.Wrapf(err, "unexpected response (%d)", response.StatusCode)
	}

	// response is good, read the bytes
	return response.Body, nil
}

func (c *proxyClient) newRequest(uri string) (*http.Request, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		uri,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create request")
	}
	return request, nil
}
