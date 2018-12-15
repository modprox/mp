package get

import (
	"archive/zip"
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/modprox/mp/pkg/clients/zips"
	"github.com/modprox/mp/pkg/clients/zips/zipstest"
	"github.com/modprox/mp/pkg/coordinates"
	"github.com/modprox/mp/pkg/metrics/stats/statstest"
	"github.com/modprox/mp/pkg/repository"
	"github.com/modprox/mp/pkg/upstream"
	"github.com/modprox/mp/pkg/upstream/upstreamtest"
	"github.com/modprox/mp/proxy/internal/modules/store"
	"github.com/modprox/mp/proxy/internal/modules/store/storetest"
)

type mocks struct {
	resolver  *upstreamtest.Resolver
	zipClient *zipstest.Client
	zipStore  *storetest.ZipStore
	index     *storetest.Index
	emitter   *statstest.Sender
}

func (m mocks) assertions(t *testing.T) {
	m.resolver.AssertExpectations(t)
	m.zipClient.AssertExpectations(t)
	m.zipStore.AssertExpectations(t)
	m.index.AssertExpectations(t)
	m.emitter.AssertExpectations(t)
}

func newMocks() mocks {
	return mocks{
		resolver:  &upstreamtest.Resolver{},
		zipClient: &zipstest.Client{},
		zipStore:  &storetest.ZipStore{},
		index:     &storetest.Index{},
		emitter:   &statstest.Sender{},
	}
}

func dummyZip(t *testing.T) repository.Blob {
	buf := new(bytes.Buffer)
	zipper := zip.NewWriter(buf)

	var files = []struct {
		Name string
		Body string
	}{
		{
			Name: "stuff/README.txt",
			Body: "this is a readme file",
		},
		{
			Name: "stuff/foo.go",
			Body: "package foo",
		},
	}

	for _, file := range files {
		f, err := zipper.Create(file.Name)
		require.NoError(t, err)
		_, err = f.Write([]byte(file.Body))
		require.NoError(t, err)
	}

	err := zipper.Close()
	require.NoError(t, err)

	return repository.Blob(buf.Bytes())
}

func Test_Download_ok(t *testing.T) {
	mocks := newMocks()
	defer mocks.assertions(t)

	serialModule := coordinates.SerialModule{
		Module: coordinates.Module{
			Source:  "github.com/pkg/errors",
			Version: "v1.2.3",
		},
		SerialID: 16,
	}

	upstreamRequest := &upstream.Request{
		Transport: "https",
		Domain:    "github.com",
		Namespace: []string{"pkg", "errors"},
		Version:   "v1.2.3",
	}

	originalBlob := dummyZip(t)
	rewrittenBlob, err := zips.Rewrite(serialModule.Module, originalBlob)
	require.NoError(t, err)

	mocks.resolver.On("Resolve", serialModule.Module).Return(
		upstreamRequest, nil,
	).Once()

	mocks.zipClient.On("Get", upstreamRequest).Return(
		originalBlob, nil,
	).Once()

	mocks.emitter.On("GaugeMS",
		"download-mod-elapsed-ms", mock.AnythingOfType("time.Time"),
	).Once()

	mocks.zipStore.On("PutZip",
		serialModule.Module,
		rewrittenBlob,
	).Return(nil).Once()

	mocks.index.On("Put", store.ModuleAddition{
		Mod:      serialModule.Module,
		UniqueID: 16,
		ModFile:  "module github.com/pkg/errors\n",
	}).Return(nil).Once()

	dl := New(
		mocks.zipClient,
		mocks.resolver,
		mocks.zipStore,
		mocks.index,
		mocks.emitter,
	)

	err = dl.Download(serialModule)
	require.NoError(t, err)
}

func Test_Download_err_Resolve(t *testing.T) {
	mocks := newMocks()
	defer mocks.assertions(t)

	serialModule := coordinates.SerialModule{
		Module: coordinates.Module{
			Source:  "github.com/pkg/errors",
			Version: "v1.2.3",
		},
		SerialID: 16,
	}

	upstreamRequest := &upstream.Request{
		Transport: "https",
		Domain:    "github.com",
		Namespace: []string{"pkg", "errors"},
		Version:   "v1.2.3",
	}

	// originalBlob := dummyZip(t)
	// rewrittenBlob, err := zips.Rewrite(serialMod.Module, originalBlob)
	// require.NoError(t, err)

	mocks.resolver.On("Resolve", serialModule.Module).Return(
		upstreamRequest, nil,
	).Once()

	mocks.zipClient.On("Get", upstreamRequest).Return(
		nil, errors.New("zip client get failed"),
	).Once()

	dl := New(
		mocks.zipClient,
		mocks.resolver,
		mocks.zipStore,
		mocks.index,
		mocks.emitter,
	)

	err := dl.Download(serialModule)
	require.EqualError(t, err, "zip client get failed")
}

func Test_Download_err_Get(t *testing.T) {
	mocks := newMocks()
	defer mocks.assertions(t)

	serialModule := coordinates.SerialModule{
		Module: coordinates.Module{
			Source:  "broken",
			Version: "broken",
		},
		SerialID: 0,
	}

	mocks.resolver.On("Resolve", serialModule.Module).Return(
		nil, errors.New("error on resolve"),
	).Once()

	dl := New(
		mocks.zipClient,
		mocks.resolver,
		mocks.zipStore,
		mocks.index,
		mocks.emitter,
	)

	err := dl.Download(serialModule)
	require.EqualError(t, err, "error on resolve")
}

func Test_Download_err_Rewrite(t *testing.T) {
	mocks := newMocks()
	defer mocks.assertions(t)

	serialModule := coordinates.SerialModule{
		Module: coordinates.Module{
			Source:  "broken",
			Version: "broken",
		},
		SerialID: 0,
	}

	upstreamRequest := &upstream.Request{
		Transport: "https",
		Domain:    "github.com",
		Namespace: []string{"pkg", "errors"},
		Version:   "v1.2.3",
	}

	// will cause zip rewrite failure (not valid zip file)
	badBlob := repository.Blob([]byte{1, 2, 3, 4})

	mocks.resolver.On("Resolve", serialModule.Module).Return(
		upstreamRequest, nil,
	).Once()

	mocks.zipClient.On("Get", upstreamRequest).Return(
		badBlob, nil,
	).Once()

	mocks.emitter.On("GaugeMS",
		"download-mod-elapsed-ms", mock.AnythingOfType("time.Time"),
	).Once()

	dl := New(
		mocks.zipClient,
		mocks.resolver,
		mocks.zipStore,
		mocks.index,
		mocks.emitter,
	)

	err := dl.Download(serialModule)
	require.EqualError(t, err, "zip: not a valid zip file")
}

func Test_Download_err_PutZip(t *testing.T) {
	mocks := newMocks()
	defer mocks.assertions(t)

	serialModule := coordinates.SerialModule{
		Module: coordinates.Module{
			Source:  "github.com/pkg/errors",
			Version: "v1.2.3",
		},
		SerialID: 16,
	}

	upstreamRequest := &upstream.Request{
		Transport: "https",
		Domain:    "github.com",
		Namespace: []string{"pkg", "errors"},
		Version:   "v1.2.3",
	}

	originalBlob := dummyZip(t)
	rewrittenBlob, err := zips.Rewrite(serialModule.Module, originalBlob)
	require.NoError(t, err)

	mocks.resolver.On("Resolve", serialModule.Module).Return(
		upstreamRequest, nil,
	).Once()

	mocks.zipClient.On("Get", upstreamRequest).Return(
		originalBlob, nil,
	).Once()

	mocks.emitter.On("GaugeMS",
		"download-mod-elapsed-ms", mock.AnythingOfType("time.Time"),
	).Once()

	mocks.zipStore.On("PutZip",
		serialModule.Module,
		rewrittenBlob,
	).Return(errors.New("put zip failure")).Once()

	dl := New(
		mocks.zipClient,
		mocks.resolver,
		mocks.zipStore,
		mocks.index,
		mocks.emitter,
	)

	err = dl.Download(serialModule)
	require.EqualError(t, err, "put zip failure")
}

func Test_Download_err_Put(t *testing.T) {
	mocks := newMocks()
	defer mocks.assertions(t)

	serialModule := coordinates.SerialModule{
		Module: coordinates.Module{
			Source:  "github.com/pkg/errors",
			Version: "v1.2.3",
		},
		SerialID: 16,
	}

	upstreamRequest := &upstream.Request{
		Transport: "https",
		Domain:    "github.com",
		Namespace: []string{"pkg", "errors"},
		Version:   "v1.2.3",
	}

	originalBlob := dummyZip(t)
	rewrittenBlob, err := zips.Rewrite(serialModule.Module, originalBlob)
	require.NoError(t, err)

	mocks.resolver.On("Resolve", serialModule.Module).Return(
		upstreamRequest, nil,
	).Once()

	mocks.zipClient.On("Get", upstreamRequest).Return(
		originalBlob, nil,
	).Once()

	mocks.emitter.On("GaugeMS",
		"download-mod-elapsed-ms", mock.AnythingOfType("time.Time"),
	).Once()

	mocks.zipStore.On("PutZip",
		serialModule.Module,
		rewrittenBlob,
	).Return(nil).Once()

	mocks.index.On("Put", store.ModuleAddition{
		Mod:      serialModule.Module,
		UniqueID: 16,
		ModFile:  "module github.com/pkg/errors\n",
	}).Return(errors.New("put failure")).Once()

	dl := New(
		mocks.zipClient,
		mocks.resolver,
		mocks.zipStore,
		mocks.index,
		mocks.emitter,
	)

	err = dl.Download(serialModule)
	require.EqualError(t, err, "put failure")
}
