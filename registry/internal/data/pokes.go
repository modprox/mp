package data

import "github.com/modprox/mp/pkg/pokes"

func (s *store) SetStartConfig(config pokes.StartConfig) error {
	_, err := s.statements[insertStartupConfigSQL].Exec(
		config.Self.Address,
		config.Self.Port,
		config.Transforms,
	)
	return err
}

func (s *store) SetHeartbeat(heartbeat pokes.Heartbeat) error {
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
