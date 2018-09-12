package heartbeat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/modprox/mp/pkg/clients/payloads"
	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/pkg/netservice"
)

const (
	heartbeatPath = "/v1/heartbeat/update"
)

type Options struct {
	Timeout    time.Duration
	Self       netservice.Instance
	Registries []netservice.Instance
}

// A Sender is used to send heartbeat status updates to the registry.
type Sender interface {
	Send(int, int) error
}

// todo: use registry.Client
type sender struct {
	httpClient *http.Client
	registries []netservice.Instance
	self       netservice.Instance
	log        loggy.Logger
}

func NewSender(options Options) Sender {
	timeout := options.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &sender{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		registries: options.Registries,
		self:       options.Self,
		log:        loggy.New("heartbeat-sender"),
	}
}

func (s *sender) Send(numPackages, numModules int) error {
	heartbeat := payloads.Heartbeat{
		Self:        s.self,
		NumPackages: numPackages,
		NumModules:  numModules,
	}

	s.log.Infof("sending a heartbeat: %s", heartbeat)

	for _, registry := range s.registries {
		err := s.trySend(registry, heartbeat)
		if err == nil { // equal
			s.log.Infof("send was successful")
			return nil
		}
		s.log.Warnf("send to %s was failed: %v", registry, err)
	}

	s.log.Errorf("unable to send heartbeat")
	return errors.New("unable to send heartbeat to any registry")
}

func (s *sender) trySend(
	registry netservice.Instance,
	heartbeat payloads.Heartbeat,
) error {

	host := fmt.Sprintf("%s:%d", registry.Address, registry.Port)
	uri := &url.URL{
		Scheme: "http",
		Host:   host,
		Path:   heartbeatPath,
	}
	theURI := uri.String()

	bs, err := json.Marshal(heartbeat)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(
		http.MethodPost,
		theURI,
		bytes.NewReader(bs),
	)

	_, err = s.httpClient.Do(request)
	return err
}
