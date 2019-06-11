package config

import (
	"fmt"
	"net/http"
	"time"

	"oss.indeed.com/go/modprox/pkg/configutil"
	"oss.indeed.com/go/modprox/pkg/metrics/stats"

	"github.com/pkg/errors"
)

type Configuration struct {
	WebServer WebServer       `json:"web_server"`
	CSRF      CSRF            `json:"csrf"`
	Database  PersistentStore `json:"database_storage"`
	Statsd    stats.Statsd    `json:"statsd"`
	Proxies   Proxies         `json:"proxies"`
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
	BindAddress   string   `json:"bind_address"`
	Port          int      `json:"port"`
	ReadTimeoutS  int      `json:"read_timeout_s"`
	WriteTimeoutS int      `json:"write_timeout_s"`
	APIKeys       []string `json:"api_keys"`
}

func (s WebServer) Server(mux http.Handler) (*http.Server, error) {
	if s.BindAddress == "" {
		return nil, errors.New("server bind address is not set")
	}

	if s.Port == 0 {
		return nil, errors.New("server port is not set")
	}

	if s.TLS.Enabled {
		if s.TLS.Certificate == "" {
			return nil, errors.New("TLS enabled, but server TLS certificate not set")
		}

		if s.TLS.Key == "" {
			return nil, errors.New("TLS enabled, but server TLS key not set")
		}
	}

	if s.ReadTimeoutS == 0 {
		s.ReadTimeoutS = 60
	}

	if s.WriteTimeoutS == 0 {
		s.WriteTimeoutS = 60
	}

	address := fmt.Sprintf("%s:%d", s.BindAddress, s.Port)
	server := &http.Server{
		Addr:         address,
		Handler:      mux,
		ReadTimeout:  seconds(s.ReadTimeoutS),
		WriteTimeout: seconds(s.WriteTimeoutS),
	}

	return server, nil
}

func seconds(s int) time.Duration {
	return time.Duration(s) * time.Second
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

func (dsn DSN) String() string {
	return fmt.Sprintf(
		"dsn:[user: %s, address: %s, database: %s, allownative: %t",
		dsn.User,
		dsn.Address,
		dsn.Database,
		dsn.AllowNativePasswords,
	)
}

type Proxies struct {
	PruneAfter int `json:"prune_after_s"`
}
