package zips

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/shoenig/toolkit"

	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/pkg/repository"
	"github.com/modprox/mp/pkg/upstream"
)

type httpClient struct {
	client  *http.Client
	options HTTPOptions
	log     loggy.Logger
}

type HTTPOptions struct {
	Timeout time.Duration
}

func NewHTTPClient(options HTTPOptions) Client {
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

func (hc *httpClient) Protocols() []string {
	return []string{"http", "https"}

}

func (hc *httpClient) Get(r *upstream.Request) (repository.Blob, error) {
	if r == nil {
		return nil, errors.New("request is nil")
	}

	uri := r.URI()
	hc.log.Tracef("making a request to %s", uri)

	request, err := hc.convert(r)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create request from %s", uri)
	}

	response, err := hc.client.Do(request)
	if err != nil {
		return nil, errors.Wrapf(err, "could not do request for %s", uri)
	}
	defer toolkit.Drain(response.Body)

	// if we get a bad response code, try to read the body and log it
	if response.StatusCode >= 400 {
		bs, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, errors.Wrapf(
				err,
				"could not read body of bad response (%d)",
				response.StatusCode,
			)
		}
		hc.log.Errorf(
			"bad response (%d), body: %s",
			response.StatusCode,
			string(bs),
		)
		return nil, errors.Wrapf(
			err,
			"unexpected response (%d)",
			response.StatusCode,
		)
	}

	// response is good, read the bytes
	return ioutil.ReadAll(response.Body)
}

func (hc *httpClient) convert(r *upstream.Request) (*http.Request, error) {
	uri := r.URI()
	request, err := http.NewRequest(
		http.MethodGet,
		uri,
		nil,
	)
	if err != nil {
		return nil, err
	}

	for k, v := range r.Headers {
		request.Header.Set(k, v)
	}

	return request, nil
}
