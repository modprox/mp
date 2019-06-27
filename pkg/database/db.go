package database

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/pkg/errors"

	"oss.indeed.com/go/modprox/pkg/setup"
)

func Connect(kind string, dsn setup.DSN) (*sql.DB, error) {
	var err error
	var db *sql.DB

	switch kind {
	case "mysql":

		db, err = connectMySQL(mysql.Config{
			Net:                  "tcp",
			User:                 dsn.User,
			Passwd:               dsn.Password,
			Addr:                 dsn.Address,
			DBName:               dsn.Database,
			AllowNativePasswords: dsn.AllowNativePasswords,
			ReadTimeout:          1 * time.Minute, // todo
			WriteTimeout:         1 * time.Minute, // todo
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to connect to mysql")
		}
	default:
		return nil, errors.Errorf("%s is not a supported database", kind)
	}

	return db, nil
}

func connectMySQL(config mysql.Config) (*sql.DB, error) {
	dsn := config.FormatDSN()
	return sql.Open("mysql", dsn)
}
