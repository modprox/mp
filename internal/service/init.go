package service

import (
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/modprox/modprox-registry/internal/repositories"
	"github.com/modprox/modprox-registry/internal/web"
	"github.com/pkg/errors"
)

type initer func(*Registry) error

func initStore(r *Registry) error {
	dsn := r.config.Index.MySQL
	db, err := repositories.Connect(mysql.Config{
		User:                 dsn.User,
		Passwd:               dsn.Password,
		Addr:                 dsn.Address,
		DBName:               dsn.Database,
		AllowNativePasswords: dsn.AllowNativePasswords,
	})
	if err != nil {
		return errors.Wrap(err, "failed to connect to mysql")
	}

	store, err := repositories.New(db)
	if err != nil {
		return errors.Wrap(err, "failed to open repository index")
	}

	r.store = store
	return nil
}

func initWebserver(r *Registry) error {
	go func(h http.Handler) {
		if err := http.ListenAndServe(":8000", h); err != nil {
			r.log.Errorf("failed to listen and serve forever: %v", err)
			panic(err)
		}
	}(web.NewRouter(r.store))
	return nil
}
