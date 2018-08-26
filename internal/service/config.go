package service

type Configuration struct {
	Index PersistentStore `json:"persistent_index"`
}

type PersistentStore struct {
	MySQL DSN `json:"mysql,omitempty"`
	// todo: add more options for storing things
}

// DSN represents the "data source name" for a database.
type DSN struct {
	User            string            `json:"user"`
	Password        string            `json:"password"`
	Address         string            `json:"address"`
	Database        string            `json:"database"`
	Parameters      map[string]string `json:"parameters"`
	ServerPublicKey string            `json:"server_public_key"`
}
