package heartbeat

import (
	"bytes"
	"encoding/json"

	"github.com/cactus/go-statsd-client/statsd"

	"github.com/modprox/mp/pkg/clients/payloads"
	"github.com/modprox/mp/pkg/clients/registry"
	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/pkg/netservice"
)

const (
	heartbeatPath = "/v1/proxy/heartbeat"
)

// A Sender is used to send heartbeat status updates to the registry.
type Sender interface {
	Send(int, int) error
}

// todo: use registry.Client
type sender struct {
	registryClient registry.Client
	self           netservice.Instance
	statter        statsd.StatSender
	log            loggy.Logger
}

func NewSender(
	self netservice.Instance,
	registryClient registry.Client,
	statter statsd.Statter,
) Sender {

	return &sender{
		registryClient: registryClient,
		self:           self,
		statter:        statter,
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
		s.statter.Inc("heartbeat-send-failure", 1, 1)
		return err
	}

	s.log.Infof("heartbeat was successfully sent!")
	s.statter.Inc("heartbeat-send-ok", 1, 1)
	return nil
}
