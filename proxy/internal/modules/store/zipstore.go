package store

import (
	"database/sql"

	"oss.indeed.com/go/modprox/pkg/config"
	"oss.indeed.com/go/modprox/pkg/coordinates"
	database "oss.indeed.com/go/modprox/pkg/db"
	"oss.indeed.com/go/modprox/pkg/loggy"
	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/pkg/repository"
)

//go:generate minimock -g -i ZipStore -s _mock.go

type ZipStore interface {
	PutZip(coordinates.Module, repository.Blob) error
	GetZip(coordinates.Module) (repository.Blob, error)
	DelZip(coordinates.Module) error
}

func Connect(kind string, dsn config.DSN, emitter stats.Sender) (*mysqlStore, error) {
	db, err := database.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return New(kind, db, emitter)
}

func New(kind string, db *sql.DB, emitter stats.Sender) (*mysqlStore, error) {
	statements, err := load(kind, db)
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
