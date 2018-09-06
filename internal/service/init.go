package service

import (
	"fmt"
	"net/http"
	"os"

	"github.com/modprox/modprox-registry/internal/data"
	"github.com/modprox/modprox-registry/internal/web"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type initer func(*Registry) error

func initStore(r *Registry) error {
	dsn := r.config.Database.MySQL
	db, err := data.Connect(mysql.Config{
		User:                 dsn.User,
		Passwd:               dsn.Password,
		Addr:                 dsn.Address,
		DBName:               dsn.Database,
		AllowNativePasswords: dsn.AllowNativePasswords,
	})
	if err != nil {
		return errors.Wrap(err, "failed to connect to mysql")
	}

	store, err := data.New(db)
	if err != nil {
		return errors.Wrap(err, "failed to open repository index")
	}

	r.store = store
	return nil
}

func initWebServer(r *Registry) error {
	router := web.NewRouter(r.store, r.config.CSRF)

	listenOn := fmt.Sprintf(
		"%s:%d",
		r.config.WebServer.BindAddress,
		r.config.WebServer.Port,
	)

	go func(h http.Handler) {
		if err := http.ListenAndServe(listenOn, h); err != nil {
			r.log.Errorf("failed to listen and serve forever")
			r.log.Errorf("caused by: %v", err)
			os.Exit(1)
		}
	}(router)

	return nil
}
