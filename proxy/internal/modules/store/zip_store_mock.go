package store

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

import (
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"oss.indeed.com/go/modprox/pkg/coordinates"
	"oss.indeed.com/go/modprox/pkg/repository"
)

// ZipStoreMock implements ZipStore
type ZipStoreMock struct {
	t minimock.Tester

	funcDelZip          func(m1 coordinates.Module) (err error)
	inspectFuncDelZip   func(m1 coordinates.Module)
	afterDelZipCounter  uint64
	beforeDelZipCounter uint64
	DelZipMock          mZipStoreMockDelZip

	funcGetZip          func(m1 coordinates.Module) (b1 repository.Blob, err error)
	inspectFuncGetZip   func(m1 coordinates.Module)
	afterGetZipCounter  uint64
	beforeGetZipCounter uint64
	GetZipMock          mZipStoreMockGetZip

	funcPutZip          func(m1 coordinates.Module, b1 repository.Blob) (err error)
	inspectFuncPutZip   func(m1 coordinates.Module, b1 repository.Blob)
	afterPutZipCounter  uint64
	beforePutZipCounter uint64
	PutZipMock          mZipStoreMockPutZip
}

// NewZipStoreMock returns a mock for ZipStore
func NewZipStoreMock(t minimock.Tester) *ZipStoreMock {
	m := &ZipStoreMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.DelZipMock = mZipStoreMockDelZip{mock: m}
	m.DelZipMock.callArgs = []*ZipStoreMockDelZipParams{}

	m.GetZipMock = mZipStoreMockGetZip{mock: m}
	m.GetZipMock.callArgs = []*ZipStoreMockGetZipParams{}

	m.PutZipMock = mZipStoreMockPutZip{mock: m}
	m.PutZipMock.callArgs = []*ZipStoreMockPutZipParams{}

	return m
}

type mZipStoreMockDelZip struct {
	mock               *ZipStoreMock
	defaultExpectation *ZipStoreMockDelZipExpectation
	expectations       []*ZipStoreMockDelZipExpectation

	callArgs []*ZipStoreMockDelZipParams
	mutex    sync.RWMutex
}

// ZipStoreMockDelZipExpectation specifies expectation struct of the ZipStore.DelZip
type ZipStoreMockDelZipExpectation struct {
	mock    *ZipStoreMock
	params  *ZipStoreMockDelZipParams
	results *ZipStoreMockDelZipResults
	Counter uint64
}

// ZipStoreMockDelZipParams contains parameters of the ZipStore.DelZip
type ZipStoreMockDelZipParams struct {
	m1 coordinates.Module
}

// ZipStoreMockDelZipResults contains results of the ZipStore.DelZip
type ZipStoreMockDelZipResults struct {
	err error
}

// Expect sets up expected params for ZipStore.DelZip
func (mmDelZip *mZipStoreMockDelZip) Expect(m1 coordinates.Module) *mZipStoreMockDelZip {
	if mmDelZip.mock.funcDelZip != nil {
		mmDelZip.mock.t.Fatalf("ZipStoreMock.DelZip mock is already set by Set")
	}

	if mmDelZip.defaultExpectation == nil {
		mmDelZip.defaultExpectation = &ZipStoreMockDelZipExpectation{}
	}

	mmDelZip.defaultExpectation.params = &ZipStoreMockDelZipParams{m1}
	for _, e := range mmDelZip.expectations {
		if minimock.Equal(e.params, mmDelZip.defaultExpectation.params) {
			mmDelZip.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmDelZip.defaultExpectation.params)
		}
	}

	return mmDelZip
}

// Inspect accepts an inspector function that has same arguments as the ZipStore.DelZip
func (mmDelZip *mZipStoreMockDelZip) Inspect(f func(m1 coordinates.Module)) *mZipStoreMockDelZip {
	if mmDelZip.mock.inspectFuncDelZip != nil {
		mmDelZip.mock.t.Fatalf("Inspect function is already set for ZipStoreMock.DelZip")
	}

	mmDelZip.mock.inspectFuncDelZip = f

	return mmDelZip
}

// Return sets up results that will be returned by ZipStore.DelZip
func (mmDelZip *mZipStoreMockDelZip) Return(err error) *ZipStoreMock {
	if mmDelZip.mock.funcDelZip != nil {
		mmDelZip.mock.t.Fatalf("ZipStoreMock.DelZip mock is already set by Set")
	}

	if mmDelZip.defaultExpectation == nil {
		mmDelZip.defaultExpectation = &ZipStoreMockDelZipExpectation{mock: mmDelZip.mock}
	}
	mmDelZip.defaultExpectation.results = &ZipStoreMockDelZipResults{err}
	return mmDelZip.mock
}

//Set uses given function f to mock the ZipStore.DelZip method
func (mmDelZip *mZipStoreMockDelZip) Set(f func(m1 coordinates.Module) (err error)) *ZipStoreMock {
	if mmDelZip.defaultExpectation != nil {
		mmDelZip.mock.t.Fatalf("Default expectation is already set for the ZipStore.DelZip method")
	}

	if len(mmDelZip.expectations) > 0 {
		mmDelZip.mock.t.Fatalf("Some expectations are already set for the ZipStore.DelZip method")
	}

	mmDelZip.mock.funcDelZip = f
	return mmDelZip.mock
}

// When sets expectation for the ZipStore.DelZip which will trigger the result defined by the following
// Then helper
func (mmDelZip *mZipStoreMockDelZip) When(m1 coordinates.Module) *ZipStoreMockDelZipExpectation {
	if mmDelZip.mock.funcDelZip != nil {
		mmDelZip.mock.t.Fatalf("ZipStoreMock.DelZip mock is already set by Set")
	}

	expectation := &ZipStoreMockDelZipExpectation{
		mock:   mmDelZip.mock,
		params: &ZipStoreMockDelZipParams{m1},
	}
	mmDelZip.expectations = append(mmDelZip.expectations, expectation)
	return expectation
}

// Then sets up ZipStore.DelZip return parameters for the expectation previously defined by the When method
func (e *ZipStoreMockDelZipExpectation) Then(err error) *ZipStoreMock {
	e.results = &ZipStoreMockDelZipResults{err}
	return e.mock
}

// DelZip implements ZipStore
func (mmDelZip *ZipStoreMock) DelZip(m1 coordinates.Module) (err error) {
	mm_atomic.AddUint64(&mmDelZip.beforeDelZipCounter, 1)
	defer mm_atomic.AddUint64(&mmDelZip.afterDelZipCounter, 1)

	if mmDelZip.inspectFuncDelZip != nil {
		mmDelZip.inspectFuncDelZip(m1)
	}

	mm_params := &ZipStoreMockDelZipParams{m1}

	// Record call args
	mmDelZip.DelZipMock.mutex.Lock()
	mmDelZip.DelZipMock.callArgs = append(mmDelZip.DelZipMock.callArgs, mm_params)
	mmDelZip.DelZipMock.mutex.Unlock()

	for _, e := range mmDelZip.DelZipMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmDelZip.DelZipMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmDelZip.DelZipMock.defaultExpectation.Counter, 1)
		mm_want := mmDelZip.DelZipMock.defaultExpectation.params
		mm_got := ZipStoreMockDelZipParams{m1}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmDelZip.t.Errorf("ZipStoreMock.DelZip got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmDelZip.DelZipMock.defaultExpectation.results
		if mm_results == nil {
			mmDelZip.t.Fatal("No results are set for the ZipStoreMock.DelZip")
		}
		return (*mm_results).err
	}
	if mmDelZip.funcDelZip != nil {
		return mmDelZip.funcDelZip(m1)
	}
	mmDelZip.t.Fatalf("Unexpected call to ZipStoreMock.DelZip. %v", m1)
	return
}

// DelZipAfterCounter returns a count of finished ZipStoreMock.DelZip invocations
func (mmDelZip *ZipStoreMock) DelZipAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmDelZip.afterDelZipCounter)
}

// DelZipBeforeCounter returns a count of ZipStoreMock.DelZip invocations
func (mmDelZip *ZipStoreMock) DelZipBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmDelZip.beforeDelZipCounter)
}

// Calls returns a list of arguments used in each call to ZipStoreMock.DelZip.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmDelZip *mZipStoreMockDelZip) Calls() []*ZipStoreMockDelZipParams {
	mmDelZip.mutex.RLock()

	argCopy := make([]*ZipStoreMockDelZipParams, len(mmDelZip.callArgs))
	copy(argCopy, mmDelZip.callArgs)

	mmDelZip.mutex.RUnlock()

	return argCopy
}

// MinimockDelZipDone returns true if the count of the DelZip invocations corresponds
// the number of defined expectations
func (m *ZipStoreMock) MinimockDelZipDone() bool {
	for _, e := range m.DelZipMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.DelZipMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterDelZipCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcDelZip != nil && mm_atomic.LoadUint64(&m.afterDelZipCounter) < 1 {
		return false
	}
	return true
}

// MinimockDelZipInspect logs each unmet expectation
func (m *ZipStoreMock) MinimockDelZipInspect() {
	for _, e := range m.DelZipMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ZipStoreMock.DelZip with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.DelZipMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterDelZipCounter) < 1 {
		if m.DelZipMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to ZipStoreMock.DelZip")
		} else {
			m.t.Errorf("Expected call to ZipStoreMock.DelZip with params: %#v", *m.DelZipMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcDelZip != nil && mm_atomic.LoadUint64(&m.afterDelZipCounter) < 1 {
		m.t.Error("Expected call to ZipStoreMock.DelZip")
	}
}

type mZipStoreMockGetZip struct {
	mock               *ZipStoreMock
	defaultExpectation *ZipStoreMockGetZipExpectation
	expectations       []*ZipStoreMockGetZipExpectation

	callArgs []*ZipStoreMockGetZipParams
	mutex    sync.RWMutex
}

// ZipStoreMockGetZipExpectation specifies expectation struct of the ZipStore.GetZip
type ZipStoreMockGetZipExpectation struct {
	mock    *ZipStoreMock
	params  *ZipStoreMockGetZipParams
	results *ZipStoreMockGetZipResults
	Counter uint64
}

// ZipStoreMockGetZipParams contains parameters of the ZipStore.GetZip
type ZipStoreMockGetZipParams struct {
	m1 coordinates.Module
}

// ZipStoreMockGetZipResults contains results of the ZipStore.GetZip
type ZipStoreMockGetZipResults struct {
	b1  repository.Blob
	err error
}

// Expect sets up expected params for ZipStore.GetZip
func (mmGetZip *mZipStoreMockGetZip) Expect(m1 coordinates.Module) *mZipStoreMockGetZip {
	if mmGetZip.mock.funcGetZip != nil {
		mmGetZip.mock.t.Fatalf("ZipStoreMock.GetZip mock is already set by Set")
	}

	if mmGetZip.defaultExpectation == nil {
		mmGetZip.defaultExpectation = &ZipStoreMockGetZipExpectation{}
	}

	mmGetZip.defaultExpectation.params = &ZipStoreMockGetZipParams{m1}
	for _, e := range mmGetZip.expectations {
		if minimock.Equal(e.params, mmGetZip.defaultExpectation.params) {
			mmGetZip.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetZip.defaultExpectation.params)
		}
	}

	return mmGetZip
}

// Inspect accepts an inspector function that has same arguments as the ZipStore.GetZip
func (mmGetZip *mZipStoreMockGetZip) Inspect(f func(m1 coordinates.Module)) *mZipStoreMockGetZip {
	if mmGetZip.mock.inspectFuncGetZip != nil {
		mmGetZip.mock.t.Fatalf("Inspect function is already set for ZipStoreMock.GetZip")
	}

	mmGetZip.mock.inspectFuncGetZip = f

	return mmGetZip
}

// Return sets up results that will be returned by ZipStore.GetZip
func (mmGetZip *mZipStoreMockGetZip) Return(b1 repository.Blob, err error) *ZipStoreMock {
	if mmGetZip.mock.funcGetZip != nil {
		mmGetZip.mock.t.Fatalf("ZipStoreMock.GetZip mock is already set by Set")
	}

	if mmGetZip.defaultExpectation == nil {
		mmGetZip.defaultExpectation = &ZipStoreMockGetZipExpectation{mock: mmGetZip.mock}
	}
	mmGetZip.defaultExpectation.results = &ZipStoreMockGetZipResults{b1, err}
	return mmGetZip.mock
}

//Set uses given function f to mock the ZipStore.GetZip method
func (mmGetZip *mZipStoreMockGetZip) Set(f func(m1 coordinates.Module) (b1 repository.Blob, err error)) *ZipStoreMock {
	if mmGetZip.defaultExpectation != nil {
		mmGetZip.mock.t.Fatalf("Default expectation is already set for the ZipStore.GetZip method")
	}

	if len(mmGetZip.expectations) > 0 {
		mmGetZip.mock.t.Fatalf("Some expectations are already set for the ZipStore.GetZip method")
	}

	mmGetZip.mock.funcGetZip = f
	return mmGetZip.mock
}

// When sets expectation for the ZipStore.GetZip which will trigger the result defined by the following
// Then helper
func (mmGetZip *mZipStoreMockGetZip) When(m1 coordinates.Module) *ZipStoreMockGetZipExpectation {
	if mmGetZip.mock.funcGetZip != nil {
		mmGetZip.mock.t.Fatalf("ZipStoreMock.GetZip mock is already set by Set")
	}

	expectation := &ZipStoreMockGetZipExpectation{
		mock:   mmGetZip.mock,
		params: &ZipStoreMockGetZipParams{m1},
	}
	mmGetZip.expectations = append(mmGetZip.expectations, expectation)
	return expectation
}

// Then sets up ZipStore.GetZip return parameters for the expectation previously defined by the When method
func (e *ZipStoreMockGetZipExpectation) Then(b1 repository.Blob, err error) *ZipStoreMock {
	e.results = &ZipStoreMockGetZipResults{b1, err}
	return e.mock
}

// GetZip implements ZipStore
func (mmGetZip *ZipStoreMock) GetZip(m1 coordinates.Module) (b1 repository.Blob, err error) {
	mm_atomic.AddUint64(&mmGetZip.beforeGetZipCounter, 1)
	defer mm_atomic.AddUint64(&mmGetZip.afterGetZipCounter, 1)

	if mmGetZip.inspectFuncGetZip != nil {
		mmGetZip.inspectFuncGetZip(m1)
	}

	mm_params := &ZipStoreMockGetZipParams{m1}

	// Record call args
	mmGetZip.GetZipMock.mutex.Lock()
	mmGetZip.GetZipMock.callArgs = append(mmGetZip.GetZipMock.callArgs, mm_params)
	mmGetZip.GetZipMock.mutex.Unlock()

	for _, e := range mmGetZip.GetZipMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.b1, e.results.err
		}
	}

	if mmGetZip.GetZipMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetZip.GetZipMock.defaultExpectation.Counter, 1)
		mm_want := mmGetZip.GetZipMock.defaultExpectation.params
		mm_got := ZipStoreMockGetZipParams{m1}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGetZip.t.Errorf("ZipStoreMock.GetZip got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGetZip.GetZipMock.defaultExpectation.results
		if mm_results == nil {
			mmGetZip.t.Fatal("No results are set for the ZipStoreMock.GetZip")
		}
		return (*mm_results).b1, (*mm_results).err
	}
	if mmGetZip.funcGetZip != nil {
		return mmGetZip.funcGetZip(m1)
	}
	mmGetZip.t.Fatalf("Unexpected call to ZipStoreMock.GetZip. %v", m1)
	return
}

// GetZipAfterCounter returns a count of finished ZipStoreMock.GetZip invocations
func (mmGetZip *ZipStoreMock) GetZipAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetZip.afterGetZipCounter)
}

// GetZipBeforeCounter returns a count of ZipStoreMock.GetZip invocations
func (mmGetZip *ZipStoreMock) GetZipBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetZip.beforeGetZipCounter)
}

// Calls returns a list of arguments used in each call to ZipStoreMock.GetZip.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetZip *mZipStoreMockGetZip) Calls() []*ZipStoreMockGetZipParams {
	mmGetZip.mutex.RLock()

	argCopy := make([]*ZipStoreMockGetZipParams, len(mmGetZip.callArgs))
	copy(argCopy, mmGetZip.callArgs)

	mmGetZip.mutex.RUnlock()

	return argCopy
}

// MinimockGetZipDone returns true if the count of the GetZip invocations corresponds
// the number of defined expectations
func (m *ZipStoreMock) MinimockGetZipDone() bool {
	for _, e := range m.GetZipMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetZipMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetZipCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetZip != nil && mm_atomic.LoadUint64(&m.afterGetZipCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetZipInspect logs each unmet expectation
func (m *ZipStoreMock) MinimockGetZipInspect() {
	for _, e := range m.GetZipMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ZipStoreMock.GetZip with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetZipMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetZipCounter) < 1 {
		if m.GetZipMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to ZipStoreMock.GetZip")
		} else {
			m.t.Errorf("Expected call to ZipStoreMock.GetZip with params: %#v", *m.GetZipMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetZip != nil && mm_atomic.LoadUint64(&m.afterGetZipCounter) < 1 {
		m.t.Error("Expected call to ZipStoreMock.GetZip")
	}
}

type mZipStoreMockPutZip struct {
	mock               *ZipStoreMock
	defaultExpectation *ZipStoreMockPutZipExpectation
	expectations       []*ZipStoreMockPutZipExpectation

	callArgs []*ZipStoreMockPutZipParams
	mutex    sync.RWMutex
}

// ZipStoreMockPutZipExpectation specifies expectation struct of the ZipStore.PutZip
type ZipStoreMockPutZipExpectation struct {
	mock    *ZipStoreMock
	params  *ZipStoreMockPutZipParams
	results *ZipStoreMockPutZipResults
	Counter uint64
}

// ZipStoreMockPutZipParams contains parameters of the ZipStore.PutZip
type ZipStoreMockPutZipParams struct {
	m1 coordinates.Module
	b1 repository.Blob
}

// ZipStoreMockPutZipResults contains results of the ZipStore.PutZip
type ZipStoreMockPutZipResults struct {
	err error
}

// Expect sets up expected params for ZipStore.PutZip
func (mmPutZip *mZipStoreMockPutZip) Expect(m1 coordinates.Module, b1 repository.Blob) *mZipStoreMockPutZip {
	if mmPutZip.mock.funcPutZip != nil {
		mmPutZip.mock.t.Fatalf("ZipStoreMock.PutZip mock is already set by Set")
	}

	if mmPutZip.defaultExpectation == nil {
		mmPutZip.defaultExpectation = &ZipStoreMockPutZipExpectation{}
	}

	mmPutZip.defaultExpectation.params = &ZipStoreMockPutZipParams{m1, b1}
	for _, e := range mmPutZip.expectations {
		if minimock.Equal(e.params, mmPutZip.defaultExpectation.params) {
			mmPutZip.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmPutZip.defaultExpectation.params)
		}
	}

	return mmPutZip
}

// Inspect accepts an inspector function that has same arguments as the ZipStore.PutZip
func (mmPutZip *mZipStoreMockPutZip) Inspect(f func(m1 coordinates.Module, b1 repository.Blob)) *mZipStoreMockPutZip {
	if mmPutZip.mock.inspectFuncPutZip != nil {
		mmPutZip.mock.t.Fatalf("Inspect function is already set for ZipStoreMock.PutZip")
	}

	mmPutZip.mock.inspectFuncPutZip = f

	return mmPutZip
}

// Return sets up results that will be returned by ZipStore.PutZip
func (mmPutZip *mZipStoreMockPutZip) Return(err error) *ZipStoreMock {
	if mmPutZip.mock.funcPutZip != nil {
		mmPutZip.mock.t.Fatalf("ZipStoreMock.PutZip mock is already set by Set")
	}

	if mmPutZip.defaultExpectation == nil {
		mmPutZip.defaultExpectation = &ZipStoreMockPutZipExpectation{mock: mmPutZip.mock}
	}
	mmPutZip.defaultExpectation.results = &ZipStoreMockPutZipResults{err}
	return mmPutZip.mock
}

//Set uses given function f to mock the ZipStore.PutZip method
func (mmPutZip *mZipStoreMockPutZip) Set(f func(m1 coordinates.Module, b1 repository.Blob) (err error)) *ZipStoreMock {
	if mmPutZip.defaultExpectation != nil {
		mmPutZip.mock.t.Fatalf("Default expectation is already set for the ZipStore.PutZip method")
	}

	if len(mmPutZip.expectations) > 0 {
		mmPutZip.mock.t.Fatalf("Some expectations are already set for the ZipStore.PutZip method")
	}

	mmPutZip.mock.funcPutZip = f
	return mmPutZip.mock
}

// When sets expectation for the ZipStore.PutZip which will trigger the result defined by the following
// Then helper
func (mmPutZip *mZipStoreMockPutZip) When(m1 coordinates.Module, b1 repository.Blob) *ZipStoreMockPutZipExpectation {
	if mmPutZip.mock.funcPutZip != nil {
		mmPutZip.mock.t.Fatalf("ZipStoreMock.PutZip mock is already set by Set")
	}

	expectation := &ZipStoreMockPutZipExpectation{
		mock:   mmPutZip.mock,
		params: &ZipStoreMockPutZipParams{m1, b1},
	}
	mmPutZip.expectations = append(mmPutZip.expectations, expectation)
	return expectation
}

// Then sets up ZipStore.PutZip return parameters for the expectation previously defined by the When method
func (e *ZipStoreMockPutZipExpectation) Then(err error) *ZipStoreMock {
	e.results = &ZipStoreMockPutZipResults{err}
	return e.mock
}

// PutZip implements ZipStore
func (mmPutZip *ZipStoreMock) PutZip(m1 coordinates.Module, b1 repository.Blob) (err error) {
	mm_atomic.AddUint64(&mmPutZip.beforePutZipCounter, 1)
	defer mm_atomic.AddUint64(&mmPutZip.afterPutZipCounter, 1)

	if mmPutZip.inspectFuncPutZip != nil {
		mmPutZip.inspectFuncPutZip(m1, b1)
	}

	mm_params := &ZipStoreMockPutZipParams{m1, b1}

	// Record call args
	mmPutZip.PutZipMock.mutex.Lock()
	mmPutZip.PutZipMock.callArgs = append(mmPutZip.PutZipMock.callArgs, mm_params)
	mmPutZip.PutZipMock.mutex.Unlock()

	for _, e := range mmPutZip.PutZipMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmPutZip.PutZipMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmPutZip.PutZipMock.defaultExpectation.Counter, 1)
		mm_want := mmPutZip.PutZipMock.defaultExpectation.params
		mm_got := ZipStoreMockPutZipParams{m1, b1}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmPutZip.t.Errorf("ZipStoreMock.PutZip got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmPutZip.PutZipMock.defaultExpectation.results
		if mm_results == nil {
			mmPutZip.t.Fatal("No results are set for the ZipStoreMock.PutZip")
		}
		return (*mm_results).err
	}
	if mmPutZip.funcPutZip != nil {
		return mmPutZip.funcPutZip(m1, b1)
	}
	mmPutZip.t.Fatalf("Unexpected call to ZipStoreMock.PutZip. %v %v", m1, b1)
	return
}

// PutZipAfterCounter returns a count of finished ZipStoreMock.PutZip invocations
func (mmPutZip *ZipStoreMock) PutZipAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmPutZip.afterPutZipCounter)
}

// PutZipBeforeCounter returns a count of ZipStoreMock.PutZip invocations
func (mmPutZip *ZipStoreMock) PutZipBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmPutZip.beforePutZipCounter)
}

// Calls returns a list of arguments used in each call to ZipStoreMock.PutZip.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmPutZip *mZipStoreMockPutZip) Calls() []*ZipStoreMockPutZipParams {
	mmPutZip.mutex.RLock()

	argCopy := make([]*ZipStoreMockPutZipParams, len(mmPutZip.callArgs))
	copy(argCopy, mmPutZip.callArgs)

	mmPutZip.mutex.RUnlock()

	return argCopy
}

// MinimockPutZipDone returns true if the count of the PutZip invocations corresponds
// the number of defined expectations
func (m *ZipStoreMock) MinimockPutZipDone() bool {
	for _, e := range m.PutZipMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.PutZipMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterPutZipCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcPutZip != nil && mm_atomic.LoadUint64(&m.afterPutZipCounter) < 1 {
		return false
	}
	return true
}

// MinimockPutZipInspect logs each unmet expectation
func (m *ZipStoreMock) MinimockPutZipInspect() {
	for _, e := range m.PutZipMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ZipStoreMock.PutZip with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.PutZipMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterPutZipCounter) < 1 {
		if m.PutZipMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to ZipStoreMock.PutZip")
		} else {
			m.t.Errorf("Expected call to ZipStoreMock.PutZip with params: %#v", *m.PutZipMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcPutZip != nil && mm_atomic.LoadUint64(&m.afterPutZipCounter) < 1 {
		m.t.Error("Expected call to ZipStoreMock.PutZip")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ZipStoreMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockDelZipInspect()

		m.MinimockGetZipInspect()

		m.MinimockPutZipInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ZipStoreMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *ZipStoreMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockDelZipDone() &&
		m.MinimockGetZipDone() &&
		m.MinimockPutZipDone()
}
