package startup

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/modprox/mp/pkg/metrics/stats"

	"github.com/stretchr/testify/require"

	"github.com/modprox/mp/pkg/clients/payloads"
	"github.com/modprox/mp/pkg/clients/registry"
	"github.com/modprox/mp/pkg/netservice"
	"github.com/modprox/mp/pkg/webutil"
	"github.com/modprox/mp/proxy/config"
)

func Test_Send_firstTry(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("some reply"))
		}),
	)
	defer ts.Close()

	address, port := webutil.ParseURL(t, ts.URL)

	emitter := stats.Discard()

	client := registry.NewClient(registry.Options{
		Timeout: 1 * time.Second,
		Instances: []netservice.Instance{{
			Address: address,
			Port:    port,
		}},
	})

	apiClient := NewSender(client, 1*time.Second, emitter)

	instance := netservice.Instance{}
	storage := config.Storage{}
	registries := config.Registry{}
	transforms := config.Transforms{}

	err := apiClient.Send(payloads.Configuration{
		Self:       instance,
		Storage:    storage,
		Registry:   registries,
		Transforms: transforms,
	})
	require.NoError(t, err)
}

func Test_Send_secondTry(t *testing.T) {
	firstTry := true
	executedSecondTry := false
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if firstTry {
				w.WriteHeader(http.StatusInternalServerError)
				firstTry = false
			} else {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("some reply"))
				executedSecondTry = true
			}
		}),
	)
	defer ts.Close()

	address, port := webutil.ParseURL(t, ts.URL)

	emitter := stats.Discard()

	client := registry.NewClient(registry.Options{
		Timeout: 1 * time.Second,
		Instances: []netservice.Instance{{
			Address: address,
			Port:    port,
		}},
	})

	apiClient := NewSender(client, 10*time.Millisecond, emitter)

	instance := netservice.Instance{}
	storage := config.Storage{}
	registries := config.Registry{}
	transforms := config.Transforms{}

	err := apiClient.Send(payloads.Configuration{
		Self:       instance,
		Storage:    storage,
		Registry:   registries,
		Transforms: transforms,
	})
	require.NoError(t, err)
	require.True(t, executedSecondTry)
}
