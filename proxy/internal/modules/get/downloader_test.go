package get

import (
	"archive/zip"
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"oss.indeed.com/go/modprox/pkg/clients/zips"
	"oss.indeed.com/go/modprox/pkg/coordinates"
	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/pkg/repository"
	"oss.indeed.com/go/modprox/pkg/upstream"
	"oss.indeed.com/go/modprox/proxy/internal/modules/store"
)

type mocks struct {
	resolver       *upstream.ResolverMock
	proxyClient    *zips.ProxyClientMock
	upstreamClient *zips.UpstreamClientMock
	zipStore       *store.ZipStoreMock
	index          *store.IndexMock
	emitter        *stats.SenderMock
}

func (m mocks) assertions() {
	m.proxyClient.MinimockFinish()
	m.upstreamClient.MinimockFinish()
	m.resolver.MinimockFinish()
	m.zipStore.MinimockFinish()
	m.index.MinimockFinish()
	m.emitter.MinimockFinish()
}

func newMocks(t *testing.T) mocks {
	return mocks{
		proxyClient:    zips.NewProxyClientMock(t),
		upstreamClient: zips.NewUpstreamClientMock(t),
		resolver:       upstream.NewResolverMock(t),
		zipStore:       store.NewZipStoreMock(t),
		index:          store.NewIndexMock(t),
		emitter:        stats.NewSenderMock(t),
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

func Test_Download_upstream_ok(t *testing.T) {
	mocks := newMocks(t)
	defer mocks.assertions()

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

	// force this module to be requested from upstream, not from global proxy
	mocks.resolver.UseProxyMock.When(serialModule.Module).Then(false, nil)

	// since we're going upstream, need to resolve the request
	mocks.resolver.ResolveMock.When(serialModule.Module).Then(upstreamRequest, nil)

	// return the original raw archive blob from upstream
	mocks.upstreamClient.GetMock.When(upstreamRequest).Then(originalBlob, nil)

	mocks.emitter.GaugeMSMock.Set(func(metric string, now time.Time) {
		require.Equal(t, "download-mod-elapsed-ms", metric)
		_ = now // ignore
	})

	mocks.zipStore.PutZipMock.When(serialModule.Module, rewrittenBlob).Then(nil)

	mocks.index.PutMock.When(store.ModuleAddition{
		Mod:      serialModule.Module,
		UniqueID: 16,
		ModFile:  "module github.com/pkg/errors\n",
	}).Then(nil)

	dl := New(
		mocks.proxyClient,
		mocks.upstreamClient,
		mocks.resolver,
		mocks.zipStore,
		mocks.index,
		mocks.emitter,
	)

	err = dl.Download(serialModule)
	require.NoError(t, err)
}

func Test_Download_proxy_ok(t *testing.T) {
	mocks := newMocks(t)
	defer mocks.assertions()

	serialModule := coordinates.SerialModule{
		Module: coordinates.Module{
			Source:  "github.com/pkg/errors",
			Version: "v1.2.3",
		},
		SerialID: 16,
	}

	// when downloading from a proxy, the blob is already a well-formed zip
	originalBlob, err := zips.Rewrite(serialModule.Module, dummyZip(t))
	require.NoError(t, err)

	// allow this module to be requested from a global proxy
	mocks.resolver.UseProxyMock.When(serialModule.Module).Then(true, nil)

	// return the well-formed zip in response
	mocks.proxyClient.GetMock.When(serialModule.Module).Then(originalBlob, nil)

	mocks.emitter.GaugeMSMock.Set(func(metric string, now time.Time) {
		require.Equal(t, "download-mod-elapsed-ms", metric)
		_ = now // ignore
	})

	mocks.zipStore.PutZipMock.When(serialModule.Module, originalBlob).Then(nil)

	mocks.index.PutMock.When(store.ModuleAddition{
		Mod:      serialModule.Module,
		UniqueID: 16,
		ModFile:  "module github.com/pkg/errors\n",
	}).Then(nil)

	dl := New(
		mocks.proxyClient,
		mocks.upstreamClient,
		mocks.resolver,
		mocks.zipStore,
		mocks.index,
		mocks.emitter,
	)

	err = dl.Download(serialModule)
	require.NoError(t, err)
}

/*
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

	//mocks.resolver.On("Resolve", serialModule.Module).Return(
	//	upstreamRequest, nil,
	//).Once()
	//
	//mocks.upstreamClient.On("Get", upstreamRequest).Return(
	//	nil, errors.New("zip client get failed"),
	//).Once()

	dl := New(
		mocks.upstreamClient,
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

	//mocks.resolver.On("Resolve", serialModule.Module).Return(
	//	nil, errors.New("error on resolve"),
	//).Once()

	dl := New(
		mocks.upstreamClient,
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

	//mocks.resolver.On("Resolve", serialModule.Module).Return(
	//	upstreamRequest, nil,
	//).Once()
	//
	//mocks.upstreamClient.On("Get", upstreamRequest).Return(
	//	badBlob, nil,
	//).Once()
	//
	//mocks.emitter.On("GaugeMS",
	//	"download-mod-elapsed-ms", mock.AnythingOfType("time.Time"),
	//).Once()

	dl := New(
		mocks.upstreamClient,
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

	mocks.upstreamClient.On("Get", upstreamRequest).Return(
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
		mocks.upstreamClient,
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

	mocks.upstreamClient.On("Get", upstreamRequest).Return(
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
		mocks.upstreamClient,
		mocks.resolver,
		mocks.zipStore,
		mocks.index,
		mocks.emitter,
	)

	err = dl.Download(serialModule)
	require.EqualError(t, err, "put failure")
}
*/
