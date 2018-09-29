package web

import (
	"net/http"
	"testing"

	"github.com/modprox/mp/registry/internal/data/datatest"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/stretchr/testify/require"
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

	mocks := newMocks()

	router := NewRouter(nil, nil, mocks.store, statter)
	return router, mocks
}
