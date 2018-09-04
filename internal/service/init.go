package service

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/csrf"
	"github.com/pkg/errors"

	"github.com/modprox/modprox-registry/internal/data"
	"github.com/modprox/modprox-registry/internal/web"
)

type initer func(*Registry) error
type middleware func(http.Handler) http.Handler

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

// chain recursively chains middleware together
func chain(h http.Handler, m ...middleware) http.Handler {
	if len(m) == 0 {
		return h
	}
	return m[0](chain(h, m[1:cap(m)]...))
}

func initWebServer(r *Registry) error {
	key, err := r.config.csrfKey()
	if err != nil {
		return err
	}

	middlewares := []middleware{
		csrf.Protect(
			// the key is used to generate csrf tokens to hand
			// out on html form loads
			key,

			// CSRF cookies are https-only normally, so for development
			// mode make sure the csrf package knows we are using http
			csrf.Secure(!r.config.CSRF.DevelopmentMode), //todo: also if no tls is set?
		),
	}

	router := chain(web.NewRouter(r.store), middlewares...)

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
