package config

import (
	"errors"
	"fmt"
)

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
