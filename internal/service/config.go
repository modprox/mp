package service

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

func (c Configuration) csrfKey() ([]byte, error) {
	key := c.CSRF.AuthenticationKey
	if len(key) != 32 {
		return nil, errors.Errorf(
			"csrf.authentication_key must be 32 bytes long, got %d",
			len(key),
		)
	}
	return []byte(key), nil
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
