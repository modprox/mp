package service

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/csrf"

	"github.com/pkg/errors"

	"go.gophers.dev/pkgs/repeat/x"

	"oss.indeed.com/go/modprox/pkg/history"
	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/pkg/webutil"
	"oss.indeed.com/go/modprox/registry/internal/data"
	"oss.indeed.com/go/modprox/registry/internal/proxies"
	"oss.indeed.com/go/modprox/registry/internal/web"
)

type initer func(*Registry) error

func initSender(r *Registry) error {
	cfg := r.config.Statsd.Agent
	if cfg.Port == 0 || cfg.Address == "" {
		r.emitter = stats.Discard()
		r.log.Warnf("stats emitter is set to discard client - no metrics will be reported")
		return nil
	}

	emitter, err := stats.New(stats.Registry, r.config.Statsd)
	if err != nil {
		return err
	}
	r.emitter = emitter
	return nil
}

func initStore(r *Registry) error {
	kind, dsn, err := r.config.Database.DSN()
	if err != nil {
		return errors.Wrap(err, "failed to configure database")
	}
	r.log.Infof("using database of kind: %q", kind)
	r.log.Infof("database dsn: %s", dsn)
	store, err := data.Connect(kind, dsn, r.emitter)
	r.store = store
	return err
}

func initProxyPrune(r *Registry) error {
	maxAge := time.Duration(r.config.Proxies.PruneAfter) * time.Second
	pruner := proxies.NewPruner(maxAge, r.store)
	go func() {
		_ = x.Interval(1*time.Minute, func() error {
			_ = pruner.Prune(time.Now())
			return nil
		})
	}()
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
		r.emitter,
		r.history,
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

func initHistory(r *Registry) error {
	historyBytes, err := history.Asset("history.txt")
	if err != nil {
		return err
	}

	r.history = string(historyBytes)
	return nil
}
