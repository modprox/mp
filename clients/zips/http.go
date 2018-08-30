package zips

import (
	"errors"
	"net/http"
	"time"

	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/libmodprox/repository"
	"github.com/modprox/libmodprox/upstream"
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
		log:     loggy.New("zips-http"),
	}
}

func (c *httpClient) Protocols() []string {
	return []string{"http", "https"}

}

func (c *httpClient) Get(r *upstream.Request) (repository.Blob, error) {
	if r == nil {
		return nil, errors.New("request is nil")
	}

	uri := r.URI()
	c.log.Tracef("making a request to %s", uri)

	return nil, nil
}
