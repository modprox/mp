package web

import (
	"net/http"

	"github.com/modprox/mp/registry/internal/data/datatest"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/modprox/mp/registry/config"
	"github.com/stretchr/testify/require"

	"testing"
)

type mocks struct {
	store *datatest.Store
}

func newMocks() mocks {
	return mocks{
		store: &datatest.Store{},
	}
}

func (m mocks) assertions(t *testing.T) {
	m.store.AssertExpectations(t)
}

func makeRouter(t *testing.T) (http.Handler, mocks) {
	statter, err := statsd.NewNoop()
	require.NoError(t, err)

	csrfConfig := config.CSRF{
		DevelopmentMode:   true,
		AuthenticationKey: "12345678901234567890123456789012",
	}

	mocks := newMocks()

	router := NewRouter(nil, nil, mocks.store, csrfConfig, statter)
	return router, mocks
}
