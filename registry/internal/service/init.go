package service

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/pkg/errors"

	"github.com/modprox/mp/registry/internal/data"
	"github.com/modprox/mp/registry/internal/web"
)

type initer func(*Registry) error

func initStatter(r *Registry) error {
	var err error
	instance := r.config.StatsEmitter
	if instance.Port == 0 || instance.Address == "" {
		r.statter, err = statsd.NewNoopClient()
		r.log.Warnf("statsd statter is set to noop client")
		return err
	}
	address := fmt.Sprintf("%s:%d", instance.Address, instance.Port)
	r.statter, err = statsd.NewClient(address, "modprox-registry")
	r.log.Infof("statsd statter is set to %s", address)
	return err
}

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
	router := web.NewRouter(r.store, r.config.CSRF, r.statter)

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
