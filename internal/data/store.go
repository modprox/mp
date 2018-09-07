package data

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/libmodprox/pokes"
	"github.com/modprox/libmodprox/repository"
)

func Connect(config mysql.Config) (*sql.DB, error) {
	dsn := config.FormatDSN()
	return sql.Open("mysql", dsn)
}

type Store interface {
	// modules
	ListMods() ([]repository.ModInfo, error)
	AddMods([]repository.ModInfo) (int, int, error)

	// remove
	ListRedirects() ([]repository.Redirect, error)
	AddRedirect(repository.Redirect) error

	// startup configs and pokes
	SetStartConfig(pokes.StartConfig) error
	SetHeartbeat(pokes.Heartbeat) error
}

func New(db *sql.DB) (Store, error) {
	statements, err := load(db)
	if err != nil {
		return nil, err
	}
	return &store{
		db:         db,
		statements: statements,
		log:        loggy.New("store"),
	}, nil
}

type store struct {
	db         *sql.DB
	statements statements
	log        loggy.Logger
}

func maybeAffectedN(result sql.Result, n int) (bool, error) {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsAffected == 0 {
		return false, nil
	}

	if rowsAffected == int64(n) {
		return true, nil
	}

	return false, errors.Errorf("expected to affect %d rows, actually affected %d", n, rowsAffected)
}
