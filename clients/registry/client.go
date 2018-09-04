package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shoenig/toolkit"

	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/libmodprox/netservice"
	"github.com/modprox/libmodprox/repository"
)

// A Client is used to make requests to any one of a group of
// modprox-registry services working together.
type Client interface {
	ModInfos() ([]repository.ModInfo, error)
}

type Options struct {
	Instances []netservice.Instance
	Timeout   time.Duration
}

type client struct {
	options    Options
	httpClient *http.Client
	log        loggy.Logger
}

func NewClient(options Options) Client {
	return &client{
		options: options,
		httpClient: &http.Client{
			Timeout: options.Timeout,
		},
		log: loggy.New("registry-client"),
	}
}

func (c *client) ModInfos() ([]repository.ModInfo, error) {
	path := "/v1/registry/sources/list"
	modInfos := make([]repository.ModInfo, 0, 100)
	err := c.get(path, &modInfos)
	return modInfos, err
}

func (c *client) get(path string, i interface{}) error {
	for _, addr := range c.options.Instances {
		if err := c.getSingle(path, addr, i); err != nil {
			c.log.Warnf("GET request failed: %v", err)
			continue // keep trying with the next address
		} else {
			// the request was a success, can stop trying now
			return nil
		}
	}
	return errors.Errorf("failed to GET from any registry: %v", c.options.Instances)
}

// maybe set this in configuration somewhere
func tweak(addr string) string {
	if !strings.HasPrefix(addr, "http") {
		return "http://" + addr
	}
	return addr
}

func (c *client) getSingle(path string, addr netservice.Instance, i interface{}) error {
	url := fmt.Sprintf("%s:%d/%s", tweak(addr.Address), addr.Port, strings.TrimPrefix(path, "/"))

	c.log.Tracef("GET single for url %q", url)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		c.log.Errorf("GET single create request had error: %v", err)
		return err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		c.log.Errorf("GET single request had error: %v", err)
		return err
	}
	defer toolkit.Drain(response.Body)

	c.log.Tracef("GET single response code: %d", response.StatusCode)

	if response.StatusCode >= 400 {
		bs, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		log.Println("failed to execute GET request, body:", string(bs))
		return errors.Errorf("bad response code: %d", response.StatusCode)
	}

	return json.NewDecoder(response.Body).Decode(i)
}
