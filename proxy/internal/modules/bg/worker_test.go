package bg

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/modprox/mp/pkg/clients/registry/registrytest"
	"github.com/modprox/mp/pkg/coordinates"
	"github.com/modprox/mp/pkg/metrics/stats/statstest"
	"github.com/modprox/mp/proxy/internal/modules/get/gettest"
	"github.com/modprox/mp/proxy/internal/modules/store"
	"github.com/modprox/mp/proxy/internal/modules/store/storetest"
	"github.com/modprox/mp/proxy/internal/problems/problemstest"
)

type mocks struct {
	registryClient *registrytest.Client
	emitter        *statstest.Sender
	dlTracker      *problemstest.Tracker
	index          *storetest.Index
	store          *storetest.ZipStore
	downloader     *gettest.Downloader
	regRequester   *gettest.RegistryAPI
}

func (m mocks) assertions(t *testing.T) {
	m.registryClient.AssertExpectations(t)
	m.emitter.AssertExpectations(t)
	m.dlTracker.AssertExpectations(t)
	m.index.AssertExpectations(t)
	m.store.AssertExpectations(t)
	m.downloader.AssertExpectations(t)
	m.regRequester.AssertExpectations(t)

}

func newMocks() mocks {
	return mocks{
		registryClient: &registrytest.Client{},
		emitter:        &statstest.Sender{},
		dlTracker:      &problemstest.Tracker{},
		index:          &storetest.Index{},
		store:          &storetest.ZipStore{},
		downloader:     &gettest.Downloader{},
		regRequester:   &gettest.RegistryAPI{},
	}
}

func Test_Worker_acquireMods_ok(t *testing.T) {
	mocks := newMocks()
	defer mocks.assertions(t)

	ranges := store.Ranges{
		[2]int64{2, 3},
		[2]int64{8, 8},
	}

	smod2 := coordinates.SerialModule{
		Module: coordinates.Module{
			Source:  "github.com/foo/bar",
			Version: "v0.0.1",
		},
		SerialID: 2,
	}

	smod3 := coordinates.SerialModule{
		Module: coordinates.Module{
			Source:  "github.com/a/b",
			Version: "v0.1.2",
		},
		SerialID: 3,
	}

	smod8 := coordinates.SerialModule{
		Module: coordinates.Module{
			Source:  "github.com/a/b",
			Version: "v0.2.0",
		},
		SerialID: 8,
	}

	smods := []coordinates.SerialModule{
		smod2, smod3, smod8,
	}

	mocks.index.On("IDs").Return(ranges, nil).Once()
	mocks.regRequester.On("ModulesNeeded", ranges).Return(
		smods, nil,
	)

	// smod2 not needed
	mocks.index.On("Contains", coordinates.Module{
		Source: "github.com/foo/bar", Version: "v0.0.1",
	}).Return(true, int64(2), nil)

	// smod3 needed
	mocks.index.On("Contains", coordinates.Module{
		Source: "github.com/a/b", Version: "v0.1.2",
	}).Return(false, int64(0), nil)
	mocks.downloader.On("Download", smod3).Return(nil)

	// smod8 not needed
	mocks.index.On("Contains", coordinates.Module{
		Source: "github.com/a/b", Version: "v0.2.0",
	}).Return(true, int64(8), nil)

	w := New(
		mocks.emitter,
		mocks.dlTracker,
		mocks.index,
		mocks.store,
		mocks.regRequester,
		mocks.downloader,
	)

	mods, err := w.(*worker).acquireMods()
	require.NoError(t, err)

	fmt.Println("mods:", mods)

}
