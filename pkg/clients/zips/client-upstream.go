package zips

import (
	"github.com/pkg/errors"

	"oss.indeed.com/go/modprox/pkg/repository"
	"oss.indeed.com/go/modprox/pkg/upstream"
)

//go:generate go run github.com/gojuno/minimock/cmd/minimock -g -i UpstreamClient -s _mock.go

// UpstreamClient is used to download .zip files from an upstream origin
// (e.g. github.com). The returned Blob is in a git archive format
// that must be unpacked and repacked in the way that Go modules are
// expected to be. This is done using Rewrite.
type UpstreamClient interface {
	Get(*upstream.Request) (repository.Blob, error)
	Protocols() []string
}

func NewUpstreamClient(clients ...UpstreamClient) UpstreamClient {
	clientForProto := make(map[string]UpstreamClient, 1)
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
	clients map[string]UpstreamClient
}

func (c *client) Get(r *upstream.Request) (repository.Blob, error) {
	impl, err := c.getClientFor(r.Transport)
	if err != nil {
		return nil, err
	}
	return impl.Get(r)
}

func (c *client) Protocols() []string {
	protocols := make([]string, 0, len(c.clients))
	for proto := range c.clients {
		protocols = append(protocols, proto)
	}
	return protocols
}

func (c *client) getClientFor(transport string) (UpstreamClient, error) {
	impl, exists := c.clients[transport]
	if !exists {
		return nil, errors.Errorf("no client that handles %q", transport)
	}
	return impl, nil
}
