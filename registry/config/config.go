package config

import (
	"github.com/modprox/libmodprox/configutil"
	"github.com/pkg/errors"
)

type Configuration struct {
	WebServer WebServer       `json:"web_server"`
	CSRF      CSRF            `json:"csrf"`
	Database  PersistentStore `json:"database_storage"`
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
	MySQL DSN `json:"mysql,omitempty"`
	// todo: add more options for storing things
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
