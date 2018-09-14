package data

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/modprox/mp/pkg/clients/payloads"
	"github.com/modprox/mp/pkg/coordinates"
	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/registry/config"
)

type Store interface {
	// modules
	ListModuleIDs() ([]int64, error)
	ListModulesByIDs(ids []int64) ([]coordinates.SerialModule, error)
	ListModulesBySource(source string) ([]coordinates.SerialModule, error)
	ListModules() ([]coordinates.SerialModule, error)
	InsertModules([]coordinates.Module) (int, error)

	// startup configs and payloads
	SetStartConfig(payloads.Configuration) error
	ListStartConfigs() ([]payloads.Configuration, error)
	SetHeartbeat(payloads.Heartbeat) error
	ListHeartbeats() ([]payloads.Heartbeat, error)
}

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
