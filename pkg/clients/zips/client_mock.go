package zips

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

import (
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock"
	"oss.indeed.com/go/modprox/pkg/repository"
	"oss.indeed.com/go/modprox/pkg/upstream"
)

// ClientMock implements Client
type ClientMock struct {
	t minimock.Tester

	funcGet          func(rp1 *upstream.Request) (b1 repository.Blob, err error)
	afterGetCounter  uint64
	beforeGetCounter uint64
	GetMock          mClientMockGet

	funcProtocols          func() (sa1 []string)
	afterProtocolsCounter  uint64
	beforeProtocolsCounter uint64
	ProtocolsMock          mClientMockProtocols
}

// NewClientMock returns a mock for Client
func NewClientMock(t minimock.Tester) *ClientMock {
	m := &ClientMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}
	m.GetMock = mClientMockGet{mock: m}
	m.ProtocolsMock = mClientMockProtocols{mock: m}

	return m
}

type mClientMockGet struct {
	mock               *ClientMock
	defaultExpectation *ClientMockGetExpectation
	expectations       []*ClientMockGetExpectation
}

// ClientMockGetExpectation specifies expectation struct of the Client.Get
type ClientMockGetExpectation struct {
	mock    *ClientMock
	params  *ClientMockGetParams
	results *ClientMockGetResults
	Counter uint64
}

// ClientMockGetParams contains parameters of the Client.Get
type ClientMockGetParams struct {
	rp1 *upstream.Request
}

// ClientMockGetResults contains results of the Client.Get
type ClientMockGetResults struct {
	b1  repository.Blob
	err error
}

// Expect sets up expected params for Client.Get
func (m *mClientMockGet) Expect(rp1 *upstream.Request) *mClientMockGet {
	if m.mock.funcGet != nil {
		m.mock.t.Fatalf("ClientMock.Get mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ClientMockGetExpectation{}
	}

	m.defaultExpectation.params = &ClientMockGetParams{rp1}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Client.Get
func (m *mClientMockGet) Return(b1 repository.Blob, err error) *ClientMock {
	if m.mock.funcGet != nil {
		m.mock.t.Fatalf("ClientMock.Get mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ClientMockGetExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &ClientMockGetResults{b1, err}
	return m.mock
}

//Set uses given function f to mock the Client.Get method
func (m *mClientMockGet) Set(f func(rp1 *upstream.Request) (b1 repository.Blob, err error)) *ClientMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Client.Get method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Client.Get method")
	}

	m.mock.funcGet = f
	return m.mock
}

// When sets expectation for the Client.Get which will trigger the result defined by the following
// Then helper
func (m *mClientMockGet) When(rp1 *upstream.Request) *ClientMockGetExpectation {
	if m.mock.funcGet != nil {
		m.mock.t.Fatalf("ClientMock.Get mock is already set by Set")
	}

	expectation := &ClientMockGetExpectation{
		mock:   m.mock,
		params: &ClientMockGetParams{rp1},
	}
	m.expectations = append(m.expectations, expectation)
	return expectation
}

// Then sets up Client.Get return parameters for the expectation previously defined by the When method
func (e *ClientMockGetExpectation) Then(b1 repository.Blob, err error) *ClientMock {
	e.results = &ClientMockGetResults{b1, err}
	return e.mock
}

// Get implements Client
func (m *ClientMock) Get(rp1 *upstream.Request) (b1 repository.Blob, err error) {
	mm_atomic.AddUint64(&m.beforeGetCounter, 1)
	defer mm_atomic.AddUint64(&m.afterGetCounter, 1)

	for _, e := range m.GetMock.expectations {
		if minimock.Equal(*e.params, ClientMockGetParams{rp1}) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.b1, e.results.err
		}
	}

	if m.GetMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&m.GetMock.defaultExpectation.Counter, 1)
		want := m.GetMock.defaultExpectation.params
		got := ClientMockGetParams{rp1}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("ClientMock.Get got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := m.GetMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the ClientMock.Get")
		}
		return (*results).b1, (*results).err
	}
	if m.funcGet != nil {
		return m.funcGet(rp1)
	}
	m.t.Fatalf("Unexpected call to ClientMock.Get. %v", rp1)
	return
}

// GetAfterCounter returns a count of finished ClientMock.Get invocations
func (m *ClientMock) GetAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&m.afterGetCounter)
}

// GetBeforeCounter returns a count of ClientMock.Get invocations
func (m *ClientMock) GetBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&m.beforeGetCounter)
}

// MinimockGetDone returns true if the count of the Get invocations corresponds
// the number of defined expectations
func (m *ClientMock) MinimockGetDone() bool {
	for _, e := range m.GetMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGet != nil && mm_atomic.LoadUint64(&m.afterGetCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetInspect logs each unmet expectation
func (m *ClientMock) MinimockGetInspect() {
	for _, e := range m.GetMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ClientMock.Get with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetCounter) < 1 {
		m.t.Errorf("Expected call to ClientMock.Get with params: %#v", *m.GetMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGet != nil && mm_atomic.LoadUint64(&m.afterGetCounter) < 1 {
		m.t.Error("Expected call to ClientMock.Get")
	}
}

type mClientMockProtocols struct {
	mock               *ClientMock
	defaultExpectation *ClientMockProtocolsExpectation
	expectations       []*ClientMockProtocolsExpectation
}

// ClientMockProtocolsExpectation specifies expectation struct of the Client.Protocols
type ClientMockProtocolsExpectation struct {
	mock *ClientMock

	results *ClientMockProtocolsResults
	Counter uint64
}

// ClientMockProtocolsResults contains results of the Client.Protocols
type ClientMockProtocolsResults struct {
	sa1 []string
}

// Expect sets up expected params for Client.Protocols
func (m *mClientMockProtocols) Expect() *mClientMockProtocols {
	if m.mock.funcProtocols != nil {
		m.mock.t.Fatalf("ClientMock.Protocols mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ClientMockProtocolsExpectation{}
	}

	return m
}

// Return sets up results that will be returned by Client.Protocols
func (m *mClientMockProtocols) Return(sa1 []string) *ClientMock {
	if m.mock.funcProtocols != nil {
		m.mock.t.Fatalf("ClientMock.Protocols mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &ClientMockProtocolsExpectation{mock: m.mock}
	}
	m.defaultExpectation.results = &ClientMockProtocolsResults{sa1}
	return m.mock
}

//Set uses given function f to mock the Client.Protocols method
func (m *mClientMockProtocols) Set(f func() (sa1 []string)) *ClientMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Client.Protocols method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Client.Protocols method")
	}

	m.mock.funcProtocols = f
	return m.mock
}

// Protocols implements Client
func (m *ClientMock) Protocols() (sa1 []string) {
	mm_atomic.AddUint64(&m.beforeProtocolsCounter, 1)
	defer mm_atomic.AddUint64(&m.afterProtocolsCounter, 1)

	if m.ProtocolsMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&m.ProtocolsMock.defaultExpectation.Counter, 1)

		results := m.ProtocolsMock.defaultExpectation.results
		if results == nil {
			m.t.Fatal("No results are set for the ClientMock.Protocols")
		}
		return (*results).sa1
	}
	if m.funcProtocols != nil {
		return m.funcProtocols()
	}
	m.t.Fatalf("Unexpected call to ClientMock.Protocols.")
	return
}

// ProtocolsAfterCounter returns a count of finished ClientMock.Protocols invocations
func (m *ClientMock) ProtocolsAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&m.afterProtocolsCounter)
}

// ProtocolsBeforeCounter returns a count of ClientMock.Protocols invocations
func (m *ClientMock) ProtocolsBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&m.beforeProtocolsCounter)
}

// MinimockProtocolsDone returns true if the count of the Protocols invocations corresponds
// the number of defined expectations
func (m *ClientMock) MinimockProtocolsDone() bool {
	for _, e := range m.ProtocolsMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ProtocolsMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterProtocolsCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcProtocols != nil && mm_atomic.LoadUint64(&m.afterProtocolsCounter) < 1 {
		return false
	}
	return true
}

// MinimockProtocolsInspect logs each unmet expectation
func (m *ClientMock) MinimockProtocolsInspect() {
	for _, e := range m.ProtocolsMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to ClientMock.Protocols")
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ProtocolsMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterProtocolsCounter) < 1 {
		m.t.Error("Expected call to ClientMock.Protocols")
	}
	// if func was set then invocations count should be greater than zero
	if m.funcProtocols != nil && mm_atomic.LoadUint64(&m.afterProtocolsCounter) < 1 {
		m.t.Error("Expected call to ClientMock.Protocols")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ClientMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockGetInspect()

		m.MinimockProtocolsInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ClientMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *ClientMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetDone() &&
		m.MinimockProtocolsDone()
}
