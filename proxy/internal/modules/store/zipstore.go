package store

import (
	"database/sql"

	"oss.indeed.com/go/modprox/pkg/coordinates"
	database "oss.indeed.com/go/modprox/pkg/db"
	"oss.indeed.com/go/modprox/pkg/loggy"
	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/pkg/repository"
	"oss.indeed.com/go/modprox/pkg/setup"
)

//go:generate minimock -g -i ZipStore -s _mock.go

type ZipStore interface {
	PutZip(coordinates.Module, repository.Blob) error
	GetZip(coordinates.Module) (repository.Blob, error)
	DelZip(coordinates.Module) error
}

func Connect(dsn setup.DSN, emitter stats.Sender) (*mysqlStore, error) {
	db, err := database.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return New(db, emitter)
}

func New(db *sql.DB, emitter stats.Sender) (*mysqlStore, error) {
	statements, err := load(db)
	if err != nil {
		return nil, err
	}
	return &mysqlStore{
		emitter:    emitter,
		db:         db,
		statements: statements,
		log:        loggy.New("store"),
	}, nil
}
