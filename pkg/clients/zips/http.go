package zips

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/shoenig/httplus/responses"

	"go.gophers.dev/pkgs/loggy"

	"oss.indeed.com/go/modprox/pkg/repository"
	"oss.indeed.com/go/modprox/pkg/upstream"
)

var maxLoggedBody = 500

type httpClient struct {
	client  *http.Client
	options HTTPOptions
	log     loggy.Logger
}

type HTTPOptions struct {
	Timeout time.Duration
}

func NewHTTPClient(options HTTPOptions) UpstreamClient {
	if options.Timeout <= 0 {
		options.Timeout = 10 * time.Minute
	}
	return &httpClient{
		options: options,
		client: &http.Client{
			Timeout: options.Timeout,
		},
		log: loggy.New("zips-http"),
	}
}

func (c *httpClient) Protocols() []string {
	return []string{"http", "https"}

}

func (c *httpClient) Get(r *upstream.Request) (repository.Blob, error) {
	if r == nil {
		return nil, errors.New("request is nil")
	}

	zipURI := r.URI()
	c.log.Tracef("making zip upstream request to %s", zipURI)

	request, err := c.newRequest(r)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create request from %s", zipURI)
	}

	response, err := c.client.Do(request)
	if err != nil {
		return nil, errors.Wrapf(err, "could not do request for %s", zipURI)
	}
	defer responses.Drain(response)

	// if we get a bad response code, try to read the body and log it
	if response.StatusCode >= 400 {
		bs, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, errors.Wrapf(err, "could not read body of bad response (%d)", response.StatusCode)
		}
		body := string(bs)
		if len(body) <= maxLoggedBody {
			c.log.Errorf(
				"bad response (%d), body: %s",
				response.StatusCode,
				body,
			)
		} else {
			c.log.Errorf(
				"bad response (%d), body: %s...",
				response.StatusCode,
				body[:maxLoggedBody],
			)
		}
		return nil, errors.Wrapf(
			err,
			"unexpected response (%d)",
			response.StatusCode,
		)
	}

	// response is good, read the bytes
	return ioutil.ReadAll(response.Body)
}

func (c *httpClient) newRequest(r *upstream.Request) (*http.Request, error) {
	uri := r.URI()
	request, err := http.NewRequest(
		http.MethodGet,
		uri,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create request")
	}

	for k, v := range r.Headers {
		request.Header.Set(k, v)
	}

	return request, nil
}
