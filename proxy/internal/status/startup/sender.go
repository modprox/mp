package startup

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/modprox/mp/pkg/clients/payloads"
	"github.com/modprox/mp/pkg/clients/registry"
	"github.com/modprox/mp/pkg/loggy"
)

const (
	configurationPath = "/v1/proxy/configuration"
)

// A Sender is used to send startup configuration state to the registry.
type Sender interface {
	Send(configuration payloads.Configuration) error
}

type sender struct {
	registryClient registry.Client
	retryInterval  time.Duration
	log            loggy.Logger
}

func NewSender(registryClient registry.Client, retryInterval time.Duration) Sender {
	return &sender{
		registryClient: registryClient,
		retryInterval:  retryInterval,
		log:            loggy.New("startup-config-sender"),
	}
}

func (s *sender) Send(configuration payloads.Configuration) error {

	// optimistically try immediately to start with
	if err := s.trySend(configuration); err == nil {
		return nil
	}

	// didn't work; keep trying every 30 seconds until it works
	ticker := time.NewTicker(s.retryInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := s.trySend(configuration); err == nil {
			break
		} else {
			s.log.Warnf("failed to contact registry; will try again in 30s")
		}
	}
	return nil
}

func (s *sender) trySend(configuration payloads.Configuration) error {
	bs, err := json.Marshal(configuration)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(bs)
	response := bytes.NewBuffer(nil)
	if err := s.registryClient.Post(configurationPath, reader, response); err != nil {
		return err
	}

	s.log.Infof("startup configuration successfully sent!")
	return nil
}
