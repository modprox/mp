package service

import "github.com/modprox/libmodprox/configutil"

type Configuration struct {
	Index PersistentStore `json:"persistent_index"`
}

func (c Configuration) String() string {
	return configutil.Format(c)
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
