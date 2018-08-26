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

	"github.com/modprox/libmodprox/netutil"
	"github.com/modprox/libmodprox/repository"
)

// A Client is used to make requests to any one of a group of
// modprox-registry services working together.
type Client interface {
	ModInfos() ([]repository.ModInfo, error)
}

type Options struct {
	Addresses []netutil.Service
	Timeout   time.Duration
}

func NewClient(options Options) Client {
	return &client{
		options: options,
	}
}

type client struct {
	options    Options
	httpClient *http.Client
}

func (c *client) ModInfos() ([]repository.ModInfo, error) {
	path := "/v1/registry/sources/list"
	modInfos := make([]repository.ModInfo, 0, 100)
	err := c.get(path, modInfos)
	return modInfos, err
}

func (c *client) get(path string, i interface{}) error {
	for _, addr := range c.options.Addresses {
		if err := c.getSingle(path, addr, i); err != nil {
			return nil
		}
	}
	return errors.Errorf("failed to GET from any registry: %v", c.options.Addresses)
}

func (c *client) getSingle(path string, addr netutil.Service, i interface{}) error {
	url := fmt.Sprintf("%s:%d/%s", addr.Address, addr.Port, strings.TrimPrefix(path, "/"))

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer toolkit.Drain(response.Body)

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
