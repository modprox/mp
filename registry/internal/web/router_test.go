package web

import (
	"net/http"
	"testing"

	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/registry/internal/data/datatest"
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
	// emitter := &statstest.Sender{}
	emitter := stats.Discard()

	mocks := newMocks()

	router := NewRouter(nil, nil, mocks.store, emitter)
	return router, mocks
}
