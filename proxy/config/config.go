package config

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/modprox/mp/pkg/configutil"
	"github.com/modprox/mp/pkg/netservice"
)

type Configuration struct {
	APIServer     APIServer  `json:"api_server"`
	Registry      Registry   `json:"registry"`
	Statsd        Statsd     `json:"statsd"`
	ModuleStorage Storage    `json:"module_storage"`
	Transforms    Transforms `json:"transforms"`
}

func (c Configuration) String() string {
	return configutil.Format(c)
}

type APIServer struct {
	TLS struct {
		Enabled     bool   `json:"enabled"`
		Certificate string `json:"certificate"`
		Key         string `json:"key"`
	} `json:"tls"`
	BindAddress   string `json:"bind_address"`
	Port          int    `json:"port"`
	ReadTimeoutS  int    `json:"read_timeout_s"`
	WriteTimeoutS int    `json:"write_timeout_s"`
}

func (s APIServer) Server(mux http.Handler) (*http.Server, error) {
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

type Storage struct {
	DataPath  string `json:"data_path"`
	IndexPath string `json:"index_path"`
	TmpPath   string `json:"tmp_path"`
}

type instances = []netservice.Instance

type Registry struct {
	Instances       instances `json:"instances"`
	PollFrequencyS  int       `json:"poll_frequency_s"`
	RequestTimeoutS int       `json:"request_timeout_s"`
	APIKey          string    `json:"api_key"`
}

type Transforms struct {
	DomainGoGet []struct {
		Domain string `json:"domain"`
	} `json:"domain_go-get,omitempty"`
	DomainRedirects []struct {
		Original     string `json:"original"`
		Substitution string `json:"substitution"`
	} `json:"domain_redirects,omitempty"`
	DomainHeaders []struct {
		Domain  string            `json:"domain"`
		Headers map[string]string `json:"headers"`
	} `json:"domain_headers,omitempty"`
	DomainPath []struct {
		Domain string `json:"domain"`
		Path   string `json:"path"`
	} `json:"domain_paths,omitempty"`
	DomainTransport []struct {
		Domain    string `json:"domain"`
		Transport string `json:"transport"`
	} `json:"domain_transports,omitempty"`
}

type Statsd struct {
	Agent netservice.Instance `json:"agent"`
}
