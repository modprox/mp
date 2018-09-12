package data

import (
	"github.com/modprox/mp/pkg/clients/payloads"
)

func (s *store) SetStartConfig(config payloads.Configuration) error {
	_, err := s.statements[insertStartupConfigSQL].Exec(
	// todo: implement
	// 		config.Self.Address,
	// 		config.Self.Port,
	// 		config.Transforms,
	)
	return err
}

func (s *store) SetHeartbeat(heartbeat payloads.Heartbeat) error {
	_, err := s.statements[insertHeartbeatSQL].Exec(
		heartbeat.Self.Address,
		heartbeat.Self.Port,
		heartbeat.NumPackages,
		heartbeat.NumModules,
		heartbeat.NumPackages, // upsert
		heartbeat.NumModules,  // upsert
	)
	return err
}
