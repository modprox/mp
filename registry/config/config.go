package config

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"oss.indeed.com/go/modprox/pkg/configutil"
	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/pkg/setup"
)

type Configuration struct {
	WebServer   WebServer             `json:"web_server"`
	CSRF        CSRF                  `json:"csrf"`
	Database    setup.PersistentStore `json:"database_storage"`
	Statsd      stats.Statsd          `json:"statsd"`
	Proxies     Proxies               `json:"proxies"`
	ProxyClient ProxyClient           `json:"proxy_client"`
}

type ProxyClient struct {
	Protocol string `json:"protocol"` // e.g. "https"
	BaseURL  string `json:"base_url"` // e.g. "proxy.golang.org"
}

func (c Configuration) String() string {
	return configutil.Format(c)
}

type WebServer struct {
	TLS struct {
		Enabled     bool   `json:"enabled"`
		Certificate string `json:"certificate"`
		Key         string `json:"key"`
	} `json:"tls"`
	BindAddress   string   `json:"bind_address"`
	Port          int      `json:"port"`
	ReadTimeoutS  int      `json:"read_timeout_s"`
	WriteTimeoutS int      `json:"write_timeout_s"`
	APIKeys       []string `json:"api_keys"`
}

func (s WebServer) Server(mux http.Handler) (*http.Server, error) {
	if s.BindAddress == "" {
		return nil, errors.New("server bind address is not set")
	}

	if s.Port == 0 {
		return nil, errors.New("server port is not set")
	}

	if s.TLS.Enabled {
		if s.TLS.Certificate == "" {
			return nil, errors.New("TLS enabled, but server TLS certificate not set")
		}

		if s.TLS.Key == "" {
			return nil, errors.New("TLS enabled, but server TLS key not set")
		}
	}

	if s.ReadTimeoutS == 0 {
		s.ReadTimeoutS = 60
	}

	if s.WriteTimeoutS == 0 {
		s.WriteTimeoutS = 60
	}

	address := fmt.Sprintf("%s:%d", s.BindAddress, s.Port)
	server := &http.Server{
		Addr:         address,
		Handler:      mux,
		ReadTimeout:  seconds(s.ReadTimeoutS),
		WriteTimeout: seconds(s.WriteTimeoutS),
	}

	return server, nil
}

func seconds(s int) time.Duration {
	return time.Duration(s) * time.Second
}

type CSRF struct {
	DevelopmentMode   bool   `json:"development_mode"`
	AuthenticationKey string `json:"authentication_key"`
}

// Key returns the configured 32 byte CSRF key, and a bool indicating
// whether development mode is enabled. If the CSRF is not well formed,
// an error is returned.
func (c CSRF) Key() ([]byte, bool, error) {
	key := c.AuthenticationKey
	if len(key) != 32 {
		return nil, false, errors.Errorf(
			"csrf.authentication_key must be 32 bytes long, got %d",
			len(key),
		)
	}
	return []byte(key), c.DevelopmentMode, nil
}

type Proxies struct {
	PruneAfter int `json:"prune_after_s"`
}
