package registry

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shoenig/toolkit"

	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/pkg/netservice"
	"github.com/modprox/mp/pkg/webutil"
)

//go:generate mockery -interface=Client -package=registrytest

// A Client is used to make requests to any one of a group of
// registry services working together.
type Client interface {
	Get(path string, rw io.Writer) error
	Post(path string, body io.Reader, rw io.Writer) error
}

type Options struct {
	Instances []netservice.Instance
	APIKey    string
	Timeout   time.Duration
}

type client struct {
	options    Options
	httpClient *http.Client
	log        loggy.Logger
}

func NewClient(options Options) Client {
	if options.Timeout <= 0 {
		// some reasonable default timeout
		options.Timeout = 1 * time.Minute
	}

	return &client{
		options: options,
		httpClient: &http.Client{
			Timeout: options.Timeout,
		},
		log: loggy.New("registry-client"),
	}
}

func (c *client) Get(path string, w io.Writer) error {
	c.log.Tracef("GET %s", path)
	return c.get(path, w)
}

func (c *client) Post(path string, requestBody io.Reader, w io.Writer) error {
	c.log.Tracef("POST %s", path)
	return c.post(path, requestBody, w)
}

func (c *client) get(path string, rw io.Writer) error {
	for _, addr := range c.options.Instances {
		if err := c.getSingle(path, addr, rw); err != nil {
			c.log.Warnf("GET request failed: %v", err)
			continue // keep trying with the next instance
		} else {
			// the request was a success, can stop trying now
			return nil
		}
	}
	return errors.Errorf("failed to GET from any registry: %v", c.options.Instances)
}

func (c *client) post(path string, requestBody io.Reader, w io.Writer) error {
	for _, addr := range c.options.Instances {
		if err := c.postSingle(path, addr, requestBody, w); err != nil {
			c.log.Warnf("POST request failed: %v", err)
			continue // keep trying with the next instance
		} else {
			// the request was a success, can stop trying now
			return nil
		}
	}
	return errors.Errorf("failed to POST to any registry: %v", c.options.Instances)
}

// maybe set this in configuration somewhere
func tweak(addr string) string {
	if !strings.HasPrefix(addr, "http") {
		return "http://" + addr
	}
	return addr
}

func (c *client) getSingle(path string, instance netservice.Instance, w io.Writer) error {
	instanceURL := formatURL(instance, path)

	c.log.Tracef("GET single for url %q", instanceURL)

	request, err := http.NewRequest(http.MethodGet, instanceURL, nil)
	if err != nil {
		c.log.Errorf("GET single create request failed: %v", err)
		return err
	}
	c.setAPIKey(request)

	response, err := c.httpClient.Do(request)
	if err != nil {
		c.log.Errorf("GET single request failed: %v", err)
		return err
	}
	defer toolkit.Drain(response.Body)

	c.log.Tracef("GET single response code: %d", response.StatusCode)

	if badCode(response.StatusCode) {
		bs, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		c.log.Errorf("failed to execute GET request, body:", string(bs))
		return errors.Errorf("bad response code: %d", response.StatusCode)
	}

	// copy the response body into the provided writer
	if _, err := io.Copy(w, response.Body); err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	return nil
}

func (c *client) postSingle(
	path string,
	instance netservice.Instance,
	requestBody io.Reader,
	w io.Writer,
) error {

	instanceURL := formatURL(instance, path)

	c.log.Tracef("POST single for url %q", instanceURL)

	request, err := http.NewRequest(http.MethodPost, instanceURL, requestBody)
	if err != nil {
		c.log.Errorf("POST single create request failed: %v", err)
		return err
	}
	c.setAPIKey(request)

	response, err := c.httpClient.Do(request)
	if err != nil {
		c.log.Errorf("POST single request failed: %v", err)
		return err
	}
	defer toolkit.Drain(response.Body)

	c.log.Tracef("POST single response code: %d", response.StatusCode)

	if badCode(response.StatusCode) {
		bs, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		c.log.Errorf("failed to execute POST request, body: %q", string(bs))
		return errors.Errorf("bad response code: %d", response.StatusCode)
	}

	// copy the response body into the provided writer
	if _, err := io.Copy(w, response.Body); err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	return nil
}

func (c *client) setAPIKey(r *http.Request) {
	r.Header.Set(webutil.HeaderAPIKey, c.options.APIKey)
}

func formatURL(instance netservice.Instance, path string) string {
	// no port in URL
	if instance.Port == 0 {
		return fmt.Sprintf("%s/%s",
			tweak(instance.Address),
			strings.TrimPrefix(path, "/"),
		)
	}

	// specific port in URL
	return fmt.Sprintf(
		"%s:%d/%s",
		tweak(instance.Address),
		instance.Port,
		strings.TrimPrefix(path, "/"),
	)
}

func badCode(code int) bool {
	return code >= 400
}
