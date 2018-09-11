package background

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/modprox/mp/pkg/clients/registry"
	"github.com/modprox/mp/pkg/coordinates"
	"github.com/modprox/mp/pkg/netservice"
	"github.com/modprox/mp/proxy/internal/modules/store/storetest"
)

func parseURL(t *testing.T, tsURL string) (string, int) {
	tsURL = strings.TrimPrefix(tsURL, "http://")
	tokens := strings.Split(tsURL, ":")
	port, err := strconv.Atoi(tokens[1])
	require.NoError(t, err)
	return tokens[0], port
}

const modsReply = ` {"serials": [{
	"id": 2,
	"source": "github.com/pkg/errors",
	"version": "v0.8.0"
}]}`

func Test_ModulesNeeded(t *testing.T) {
	index := &storetest.Index{}

	ids := coordinates.RangeIDs{
		coordinates.RangeID{1, 3},
		coordinates.RangeID{6, 6},
		coordinates.RangeID{10, 20},
	}

	index.On("IDs").Return(ids, nil)

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(modsReply))
		}),
	)
	defer ts.Close()

	address, port := parseURL(t, ts.URL)
	client := registry.NewClient(registry.Options{
		Timeout: 10 * time.Second,
		Instances: []netservice.Instance{{
			Address: address,
			Port:    port,
		}},
	})

	apiClient := NewRegistryAPI(client, index)

	serialModules, err := apiClient.ModulesNeeded(ids)
	require.NoError(t, err)

	require.Equal(t, []coordinates.SerialModule{
		{
			SerialID: 2,
			Module: coordinates.Module{
				Source:  "github.com/pkg/errors",
				Version: "v0.8.0",
			},
		},
	}, serialModules)
}
