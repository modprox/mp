package data

import (
	"encoding/json"
	"time"

	"oss.indeed.com/go/modprox/pkg/clients/payloads"
	"oss.indeed.com/go/modprox/pkg/netservice"
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
		storageText,    // on dup
		registriesText, // on dup
		transformsText, // on dup
	)
	return err
}

func (s *store) ListStartConfigs() ([]payloads.Configuration, error) {
	start := time.Now()
	configs, err := s.listStartConfigs()
	if err != nil {
		s.emitter.Count("db-list-start-configs-error", 1)
		return nil, err
	}

	s.emitter.GaugeMS("db-list-start-configs-elapsed-ms", start)
	return configs, nil
}

func (s *store) listStartConfigs() ([]payloads.Configuration, error) {
	rows, err := s.statements[selectStartupConfigsSQL].Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []payloads.Configuration
	for rows.Next() {
		var (
			hostname       string
			port           int
			storageText    string
			registryText   string
			transformsText string
		)
		if err := rows.Scan(
			&hostname,
			&port,
			&storageText,
			&registryText,
			&transformsText,
		); err != nil {
			return nil, err
		}

		c, err := newConfig(
			hostname,
			port,
			storageText,
			registryText,
			transformsText,
		)
		if err != nil {
			return nil, err
		}
		configs = append(configs, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return configs, nil
}

func newConfig(
	hostname string,
	port int,
	storageText,
	registryText,
	transformsText string,
) (payloads.Configuration, error) {
	c := payloads.Configuration{
		Self: netservice.Instance{
			Address: hostname,
			Port:    port,
		},
	}

	if err := json.Unmarshal([]byte(storageText), &c.DiskStorage); err != nil {
		return c, err
	}

	if err := json.Unmarshal([]byte(registryText), &c.Registry); err != nil {
		return c, err
	}

	if err := json.Unmarshal([]byte(transformsText), &c.Transforms); err != nil {
		return c, err
	}

	return c, nil
}

func (s *store) SetHeartbeat(heartbeat payloads.Heartbeat) error {
	start := time.Now()
	err := s.setHeartbeat(heartbeat)
	if err != nil {
		s.emitter.Count("db-set-heartbeat-failure", 1)
		return err
	}

	s.emitter.GaugeMS("db-set-heartbeat-elapsed-ms", start)
	return nil
}

func (s *store) setHeartbeat(heartbeat payloads.Heartbeat) error {
	_, err := s.statements[insertHeartbeatSQL].Exec(
		heartbeat.Self.Address,
		heartbeat.Self.Port,
		heartbeat.NumModules,
		heartbeat.NumVersions,
		heartbeat.NumModules,
		heartbeat.NumVersions,
	)
	return err
}

func (s *store) ListHeartbeats() ([]payloads.Heartbeat, error) {
	start := time.Now()
	heartbeats, err := s.listHeartbeats()
	if err != nil {
		s.emitter.Count("db-list-heartbeats-failure", 1)
		return nil, err
	}

	s.emitter.GaugeMS("db-list-heartbeats-elapsed-ms", start)
	return heartbeats, nil
}

func (s *store) listHeartbeats() ([]payloads.Heartbeat, error) {
	rows, err := s.statements[selectHeartbeatsSQL].Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var heartbeats []payloads.Heartbeat
	for rows.Next() {
		var heartbeat payloads.Heartbeat
		if err := rows.Scan(
			&heartbeat.Self.Address,
			&heartbeat.Self.Port,
			&heartbeat.NumModules,
			&heartbeat.NumVersions,
			&heartbeat.Timestamp,
		); err != nil {
			return nil, err
		}
		heartbeats = append(heartbeats, heartbeat)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return heartbeats, nil
}

func (s *store) PurgeProxy(instance netservice.Instance) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1) remove heartbeat for proxy
	if _, err := tx.Stmt(s.statements[deleteHeartbeatSQL]).Exec(
		instance.Address,
		instance.Port,
	); err != nil {
		return err
	}

	// 2) remove startup configuration for proxy
	if _, err := tx.Stmt(s.statements[deleteStartupConfigSQL]).Exec(
		instance.Address,
		instance.Port,
	); err != nil {
		return err
	}

	return tx.Commit()
}
