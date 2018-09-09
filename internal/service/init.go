package service

import (
	"fmt"
	"net/http"
	"os"

	"github.com/modprox/modprox-registry/internal/data"

	"github.com/modprox/modprox-registry/internal/web"
	"github.com/pkg/errors"
)

type initer func(*Registry) error

func initStore(r *Registry) error {
	kind, dsn, err := r.config.Database.DSN()
	if err != nil {
		return errors.Wrap(err, "failed to configure database")
	}
	r.log.Infof("using database of kind: %q", kind)
	store, err := data.Connect(kind, dsn)
	r.store = store
	return err
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
