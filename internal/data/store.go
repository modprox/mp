package data

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/libmodprox/pokes"
	"github.com/modprox/libmodprox/repository"
	"github.com/modprox/modprox-registry/registry/config"
)

func Connect(kind string, dsn config.DSN) (Store, error) {
	var db *sql.DB
	var err error

	if kind == "mysql" {
		db, err = connectMySQL(mysql.Config{
			User:                 dsn.User,
			Passwd:               dsn.Password,
			Addr:                 dsn.Address,
			DBName:               dsn.Database,
			AllowNativePasswords: dsn.AllowNativePasswords,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to connect to mysql")
		}
	} else if kind == "postgres" {
		db, err = connectPostgreSQL(dsn)
		if err != nil {
			return nil, errors.Wrap(err, "failed to connect to postgres")
		}
	} else {
		return nil, errors.Errorf("%s is not a supported database", kind)
	}

	return New(kind, db)
}

func connectMySQL(config mysql.Config) (*sql.DB, error) {
	dsn := config.FormatDSN()
	return sql.Open("mysql", dsn)
}

func connectPostgreSQL(dsn config.DSN) (*sql.DB, error) {
	// "postgres://bob:secret@1.2.3.4:5432/mydb?sslmode=verify-full"
	connectStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable", // todo: enable ssl
		dsn.User,
		dsn.Password,
		dsn.Address,
		dsn.Database,
	)
	return sql.Open("postgres", connectStr)
}

type Store interface {
	// modules
	ListMods() ([]repository.ModInfo, error)
	AddMods([]repository.ModInfo) (int, error)

	// remove
	ListRedirects() ([]repository.Redirect, error)
	AddRedirect(repository.Redirect) error

	// startup configs and pokes
	SetStartConfig(pokes.StartConfig) error
	SetHeartbeat(pokes.Heartbeat) error
}

func New(kind string, db *sql.DB) (Store, error) {
	statements, err := load(kind, db)
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
