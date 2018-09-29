package service

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/shoenig/toolkit"

	"github.com/modprox/mp/registry/internal/proxies"

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
	r.log.Infof("database dsn: %s", dsn)
	store, err := data.Connect(kind, dsn, r.statter)
	r.store = store
	return err
}

func initProxyPrune(r *Registry) error {
	maxAge := time.Duration(r.config.Proxies.PruneAfter) * time.Second
	pruner := proxies.NewPruner(maxAge, r.store)
	go toolkit.Interval(1*time.Minute, func() error {
		_ = pruner.Prune(time.Now())
		return nil
	})
	return nil
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

	mux := web.NewRouter(
		middleAPI,
		middleUI,
		r.store,
		r.statter,
	)

	server, err := r.config.WebServer.Server(mux)
	if err != nil {
		return err
	}

	go func(h http.Handler) {
		var err error
		if r.config.WebServer.TLS.Enabled {
			err = server.ListenAndServeTLS(
				r.config.WebServer.TLS.Certificate,
				r.config.WebServer.TLS.Key,
			)
		} else {
			err = server.ListenAndServe()
		}

		// should never get to this point
		r.log.Errorf("server stopped serving: %v", err)
		os.Exit(1)
	}(mux)

	return nil
}
