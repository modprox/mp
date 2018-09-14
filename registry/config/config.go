package config

import (
	"github.com/modprox/mp/pkg/configutil"
	"github.com/modprox/mp/pkg/netservice"
	"github.com/pkg/errors"
)

type Configuration struct {
	WebServer    WebServer       `json:"web_server"`
	CSRF         CSRF            `json:"csrf"`
	Database     PersistentStore `json:"database_storage"`
	StatsEmitter Statsd          `json:"statsd_emitter"`
}

func (c Configuration) String() string {
	return configutil.Format(c)
}

type WebServer struct {
	TLS struct {
		Enabled     bool   `json:"enabled"`
		Certificate string `json:"certificate"`
		Key         string `json:"key"`
	} `json:"tls"`
	BindAddress string `json:"bind_address"`
	Port        int    `json:"port"`
}

type CSRF struct {
	DevelopmentMode   bool   `json:"development_mode"`
	AuthenticationKey string `json:"authentication_key"`
}

// Key returns the configured 32 byte CSRF key, and a bool indicating
// whether development mode is enabled. If the CSRF is not well formed,
// an error is returned.
func (c CSRF) Key() ([]byte, bool, error) {
	key := c.AuthenticationKey
	if len(key) != 32 {
		return nil, false, errors.Errorf(
			"csrf.authentication_key must be 32 bytes long, got %d",
			len(key),
		)
	}
	return []byte(key), c.DevelopmentMode, nil
}

type PersistentStore struct {
	MySQL      DSN `json:"mysql,omitempty"`
	PostgreSQL DSN `json:"postgres,omitempty"`
}

// DSN returns the one DSN that is configured, or returns
// an error if both or no DSN is configured.
func (ps PersistentStore) DSN() (string, DSN, error) {
	emptyDSN := DSN{}

	// check if both DSN are empty
	if ps.MySQL.equal(emptyDSN) && ps.PostgreSQL.equal(emptyDSN) {
		return "", emptyDSN, errors.New("neither mysql or postgres was configured")
	}

	// check if both DSN are configured
	if !ps.MySQL.equal(emptyDSN) && !ps.PostgreSQL.equal(emptyDSN) {
		return "", emptyDSN, errors.New("only one of mysql or postgres may be configured")
	}

	if !ps.MySQL.equal(emptyDSN) {
		return "mysql", ps.MySQL, nil
	}

	return "postgres", ps.PostgreSQL, nil
}

// DSN represents the "data source name" for a database.
type DSN struct {
	User                 string            `json:"user,omitempty"`
	Password             string            `json:"password,omitempty"`
	Address              string            `json:"address,omitempty"`
	Database             string            `json:"database,omitempty"`
	Parameters           map[string]string `json:"parameters,omitempty"`
	ServerPublicKey      string            `json:"server_public_key,omitempty"`
	AllowNativePasswords bool              `json:"allow_native_passwords,omitempty"`
}

func (dsn DSN) equal(other DSN) bool {
	switch {
	case dsn.User != other.User:
		return false
	case dsn.Password != other.Password:
		return false
	case dsn.Database != other.Database:
		return false
	case dsn.Address != other.Address:
		return false
	}
	return true
}

type Statsd netservice.Instance
