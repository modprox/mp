package data

import (
	"github.com/modprox/mp/pkg/clients/payloads"
)

func (s *store) SetStartConfig(config payloads.Configuration) error {
	storageText, registriesText, transformsText, err := config.Texts()
	if err != nil {
		return err
	}
	_, err = s.statements[insertStartupConfigSQL].Exec(
		config.Self.Address,
		config.Self.Port,
		storageText,
		registriesText,
		transformsText,
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
