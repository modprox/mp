package data

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"oss.indeed.com/go/modprox/pkg/clients/payloads"
	"oss.indeed.com/go/modprox/pkg/coordinates"
	"oss.indeed.com/go/modprox/pkg/loggy"
	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/pkg/netservice"
	"oss.indeed.com/go/modprox/registry/config"
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

func Connect(kind string, dsn config.DSN, emitter stats.Sender) (Store, error) {
	var db *sql.DB
	var err error

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
	case "postgres":
		return nil, errors.New("postgres is not supported (issue #103)")
		//db, err = connectPostgreSQL(dsn)
		//if err != nil {
		//	return nil, errors.Wrap(err, "failed to connect to postgres")
		//}
	default:
		return nil, errors.Errorf("%s is not a supported database", kind)
	}

	return New(kind, db, emitter)
}

func connectMySQL(config mysql.Config) (*sql.DB, error) {
	dsn := config.FormatDSN()
	return sql.Open("mysql", dsn)
}

//func connectPostgreSQL(dsn config.DSN) (*sql.DB, error) {
//	// "postgres://bob:secret@1.2.3.4:5432/mydb?sslmode=verify-full"
//	connectStr := fmt.Sprintf(
//		"postgres://%s:%s@%s/%s?sslmode=disable", // todo: enable ssl
//		dsn.User,
//		dsn.Password,
//		dsn.Address,
//		dsn.Database,
//	)
//	return sql.Open("postgres", connectStr)
//}

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
