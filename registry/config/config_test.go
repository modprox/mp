package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Configuration_CSRF_nokey(t *testing.T) {
	var c Configuration
	_, _, err := c.CSRF.Key()
	require.Error(t, err)
}

func Test_Configuration_CSRF_badkey(t *testing.T) {
	var c Configuration
	c.CSRF.AuthenticationKey = "foobar"
	_, _, err := c.CSRF.Key()
	require.Error(t, err)
}

func Test_Configuration_CSRF_devmode(t *testing.T) {
	key := "12345678901234567890123456789012"
	var c Configuration
	c.CSRF.DevelopmentMode = true
	c.CSRF.AuthenticationKey = key
	bs, devMode, err := c.CSRF.Key()
	require.NoError(t, err)
	require.True(t, devMode)
	require.Equal(t, []byte(key), bs)
}

func Test_Configuration_CSRF(t *testing.T) {
	key := "12345678901234567890123456789012"
	var c Configuration
	c.CSRF.AuthenticationKey = key
	bs, devMode, err := c.CSRF.Key()
	require.NoError(t, err)
	require.False(t, devMode)
	require.Equal(t, []byte(key), bs)
}

func Test_Configuration_Server_notls_defaults_ok(t *testing.T) {
	c := Configuration{
		WebServer: WebServer{
			BindAddress: "1.2.3.4",
			Port:        9999,
		},
	}

	server, err := c.WebServer.Server(nil)
	require.NoError(t, err)
	require.Equal(t, "1.2.3.4:9999", server.Addr)
	require.Equal(t, 60*time.Second, server.ReadTimeout)
	require.Equal(t, 60*time.Second, server.WriteTimeout)
}

func Test_Configuration_Server_notls_ok(t *testing.T) {
	c := Configuration{
		WebServer: WebServer{
			BindAddress:   "1.2.3.4",
			Port:          9999,
			ReadTimeoutS:  1,
			WriteTimeoutS: 2,
		},
	}

	server, err := c.WebServer.Server(nil)
	require.NoError(t, err)
	require.Equal(t, "1.2.3.4:9999", server.Addr)
	require.Equal(t, 1*time.Second, server.ReadTimeout)
	require.Equal(t, 2*time.Second, server.WriteTimeout)
}

func Test_Configuration_Server_tls_noCertificate(t *testing.T) {
	c := Configuration{
		WebServer: WebServer{
			BindAddress:   "1.2.3.4",
			Port:          9999,
			ReadTimeoutS:  1,
			WriteTimeoutS: 2,
		},
	}
	c.WebServer.TLS.Enabled = true
	c.WebServer.TLS.Key = "key"

	_, err := c.WebServer.Server(nil)
	require.Error(t, err)
}

func Test_Configuration_Server_tls_noKey(t *testing.T) {
	c := Configuration{
		WebServer: WebServer{
			BindAddress:   "1.2.3.4",
			Port:          9999,
			ReadTimeoutS:  1,
			WriteTimeoutS: 2,
		},
	}
	c.WebServer.TLS.Enabled = true
	c.WebServer.TLS.Certificate = "cert"

	_, err := c.WebServer.Server(nil)
	require.Error(t, err)
}

func Test_Configuration_Server_tls_no_files(t *testing.T) {
	c := Configuration{
		WebServer: WebServer{
			BindAddress:   "1.2.3.4",
			Port:          9999,
			ReadTimeoutS:  1,
			WriteTimeoutS: 2,
		},
	}
	c.WebServer.TLS.Enabled = true
	c.WebServer.TLS.Certificate = "cert"
	c.WebServer.TLS.Key = "key"

	server, err := c.WebServer.Server(nil)
	require.NoError(t, err)
	require.Equal(t, "1.2.3.4:9999", server.Addr)
	require.Equal(t, 1*time.Second, server.ReadTimeout)
	require.Equal(t, 2*time.Second, server.WriteTimeout)
}
