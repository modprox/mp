package web

import (
	"net/http"
	"testing"

	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/registry/internal/data"
)

type mocks struct {
	store *data.StoreMock
}

func newMocks(t *testing.T) mocks {
	return mocks{
		store: data.NewStoreMock(t),
	}
}

func (m mocks) assertions() {
	m.store.MinimockFinish()
}

func makeRouter(t *testing.T) (http.Handler, mocks) {
	// emitter := &statstest.Sender{} no testing this?

	emitter := stats.Discard()

	mocks := newMocks(t)

	router := NewRouter(nil, nil, mocks.store, emitter, "this is some fake history")
	return router, mocks
}
