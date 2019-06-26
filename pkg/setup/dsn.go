package setup

import (
	"errors"
	"fmt"
)

type PersistentStore struct {
	MySQL DSN `json:"mysql,omitempty"`
}

// DSN returns the one DSN that is configured, or returns
// an error if both or no DSN is configured.
func (ps PersistentStore) DSN() (string, DSN, error) {
	emptyDSN := DSN{}

	// check if DSN is empty
	if ps.MySQL.equal(emptyDSN) {
		return "", emptyDSN, errors.New("mysql was not configured")
	}

	return "mysql", ps.MySQL, nil
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
