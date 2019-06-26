package data

import (
	"database/sql"

	"oss.indeed.com/go/modprox/pkg/clients/payloads"
	"oss.indeed.com/go/modprox/pkg/coordinates"
	database "oss.indeed.com/go/modprox/pkg/db"
	"oss.indeed.com/go/modprox/pkg/loggy"
	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/pkg/netservice"
	"oss.indeed.com/go/modprox/pkg/setup"
)

//go:generate minimock -g -i Store -s _mock.go

type Store interface {
	// modules
	ListModuleIDs() ([]int64, error)
	ListModulesByIDs(ids []int64) ([]coordinates.SerialModule, error)
	ListModulesBySource(source string) ([]coordinates.SerialModule, error)
	ListModules() ([]coordinates.SerialModule, error)
	InsertModules([]coordinates.Module) (int, error)
	DeleteModuleByID(id int) error

	// startup configs and payloads
	SetStartConfig(payloads.Configuration) error
	ListStartConfigs() ([]payloads.Configuration, error)
	SetHeartbeat(payloads.Heartbeat) error
	ListHeartbeats() ([]payloads.Heartbeat, error)
	PurgeProxy(instance netservice.Instance) error
}

func Connect(kind string, dsn setup.DSN, emitter stats.Sender) (Store, error) {
	db, err := database.Connect(kind, dsn)
	if err != nil {
		return nil, err
	}

	return New(kind, db, emitter)
}

func New(kind string, db *sql.DB, emitter stats.Sender) (Store, error) {
	statements, err := load(kind, db)
	if err != nil {
		return nil, err
	}
	return &store{
		emitter:    emitter,
		db:         db,
		statements: statements,
		log:        loggy.New("store"),
	}, nil
}

type store struct {
	emitter    stats.Sender
	db         *sql.DB
	statements statements
	log        loggy.Logger
}
