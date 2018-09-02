package service

import (
	"fmt"

	"github.com/modprox/libmodprox/configutil"
)

type Configuration struct {
	// DevMode indicates if the server is in local development mode. Defaults to false.
	DevMode     bool            `json:"dev_mode"`
	CSRFAuthKey string          `json:"csrf_auth_key"`
	Index       PersistentStore `json:"persistent_index"`
}

func (c Configuration) String() string {
	return configutil.Format(c)
}

// MustCSRFAuthKey returns a valid 32-byte slice from configuration or panics.
func (c Configuration) MustCSRFAuthKey() []byte {
	if len(c.CSRFAuthKey) != 32 {
		panic(fmt.Sprintf("csrf_token must be 32 bytes, was: %q" + c.CSRFAuthKey))
	}
	return []byte(c.CSRFAuthKey)
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
