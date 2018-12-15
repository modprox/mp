package zips

import (
	"github.com/pkg/errors"

	"github.com/modprox/mp/pkg/repository"
	"github.com/modprox/mp/pkg/upstream"
)

//go:generate mockery3 -interface Client -package zipstest

// Client is used to download .zip files from an upstream origin
// (e.g. github.com).
type Client interface {
	Get(*upstream.Request) (repository.Blob, error)
	Protocols() []string
}

func NewClient(clients ...Client) Client {
	clientForProto := make(map[string]Client, 1)
	for _, clientImpl := range clients {
		for _, protocol := range clientImpl.Protocols() {
			clientForProto[protocol] = clientImpl
		}
	}
	return &client{
		clients: clientForProto,
	}
}

type client struct {
	clients map[string]Client
}

func (c *client) Get(r *upstream.Request) (repository.Blob, error) {
	impl, err := c.getClientFor(r.Transport)
	if err != nil {
		return nil, err
	}
	return impl.Get(r)
}

func (c *client) Protocols() []string {
	protos := make([]string, 0, len(c.clients))
	for proto := range c.clients {
		protos = append(protos, proto)
	}
	return protos
}

func (c *client) getClientFor(transport string) (Client, error) {
	impl, exists := c.clients[transport]
	if !exists {
		return nil, errors.Errorf("no client that handles %q", transport)
	}
	return impl, nil
}
