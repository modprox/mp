package service

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/gorilla/csrf"
	"github.com/pkg/errors"

	"github.com/modprox/mp/pkg/webutil"
	"github.com/modprox/mp/registry/internal/data"
	"github.com/modprox/mp/registry/internal/web"
)

type initer func(*Registry) error

func initStatter(r *Registry) error {
	var err error
	agent := r.config.Statsd.Agent
	if agent.Port == 0 || agent.Address == "" {
		r.statter, err = statsd.NewNoopClient()
		r.log.Warnf("statsd statter is set to noop client")
		return err
	}
	address := fmt.Sprintf("%s:%d", agent.Address, agent.Port)
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
	var middleAPI []webutil.Middleware
	if len(r.config.WebServer.APIKeys) > 0 {
		middleAPI = append(
			middleAPI,
			webutil.KeyGuard(r.config.WebServer.APIKeys),
		)
	}

	middleUI := []webutil.Middleware{
		csrf.Protect(
			// the key is used to generate CSRF tokens to hand
			// out on html form loads
			[]byte(r.config.CSRF.AuthenticationKey),

			// CSRF cookies are https-only normally, so for development
			//// mode make sure the CSRF package knows we are using http
			csrf.Secure(!r.config.CSRF.DevelopmentMode),
		),
	}

	router := web.NewRouter(
		middleAPI,
		middleUI,
		r.store,
		r.config.CSRF,
		r.statter,
	)

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
