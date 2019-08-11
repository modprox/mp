package store

import (
	"database/sql"

	"gophers.dev/pkgs/loggy"

	"oss.indeed.com/go/modprox/pkg/coordinates"
	"oss.indeed.com/go/modprox/pkg/database"
	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/pkg/repository"
	"oss.indeed.com/go/modprox/pkg/setup"
)

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i ZipStore -s _mock.go

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
