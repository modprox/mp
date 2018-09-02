package service

import (
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/csrf"
	"github.com/pkg/errors"

	"github.com/modprox/modprox-registry/internal/data"
	"github.com/modprox/modprox-registry/internal/web"
)

type initer func(*Registry) error
type middleware func(http.Handler) http.Handler

func initStore(r *Registry) error {
	dsn := r.config.Index.MySQL
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

func initWebserver(r *Registry) error {
	middlewares := []middleware{
		csrf.Protect(
			r.config.MustCSRFAuthKey(),
			csrf.Secure(!r.config.DevMode), // CSRF cookies are https-only normally
		),
	}
	rtr := chain(web.NewRouter(r.store), middlewares...)

	go func(h http.Handler) {
		if err := http.ListenAndServe(":8000", h); err != nil {
			r.log.Errorf("failed to listen and serve forever: %v", err)
			panic(err)
		}
	}(rtr)
	return nil
}
