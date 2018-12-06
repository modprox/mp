package heartbeat

import (
	"bytes"
	"encoding/json"

	"github.com/modprox/mp/pkg/clients/payloads"
	"github.com/modprox/mp/pkg/clients/registry"
	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/pkg/metrics/stats"
	"github.com/modprox/mp/pkg/netservice"
)

const (
	heartbeatPath = "/v1/proxy/heartbeat"
)

// A Sender is used to send heartbeat status updates to the registry.
type Sender interface {
	Send(int, int) error
}

type sender struct {
	registryClient registry.Client
	self           netservice.Instance
	emitter        stats.Sender
	log            loggy.Logger
}

func NewSender(
	self netservice.Instance,
	registryClient registry.Client,
	emitter stats.Sender,
) Sender {

	return &sender{
		registryClient: registryClient,
		self:           self,
		emitter:        emitter,
		log:            loggy.New("heartbeat-sender"),
	}
}

func (s *sender) Send(numPackages, numModules int) error {
	heartbeat := payloads.Heartbeat{
		Self:        s.self,
		NumModules:  numPackages,
		NumVersions: numModules,
	}

	s.log.Infof("sending a heartbeat: %s", heartbeat)

	bs, err := json.Marshal(heartbeat)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(bs)
	response := bytes.NewBuffer(nil)

	if err := s.registryClient.Post(heartbeatPath, reader, response); err != nil {
		s.emitter.Count("heartbeat-send-failure", 1)
		return err
	}

	s.log.Infof("heartbeat was successfully sent!")
	s.emitter.Count("heartbeat-send-ok", 1)

	return nil
}
