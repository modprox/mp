// Code autogenerated by mockery v2.0.0
//
// Do not manually edit the content of this file.

// Package storetest contains autogenerated mocks.
package storetest

import "github.com/modprox/libmodprox/coordinates"
import "github.com/stretchr/testify/mock"
import "github.com/modprox/libmodprox/repository"

// ZipStore is an autogenerated mock type for the ZipStore type
type ZipStore struct {
	mock.Mock
}

// GetZip provides a mock function with given fields: mockeryArg0
func (mockerySelf *ZipStore) GetZip(mockeryArg0 coordinates.Module) (repository.Blob, error) {
	ret := mockerySelf.Called(mockeryArg0)

	var r0 repository.Blob
	if rf, ok := ret.Get(0).(func(coordinates.Module) repository.Blob); ok {
		r0 = rf(mockeryArg0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repository.Blob)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(coordinates.Module) error); ok {
		r1 = rf(mockeryArg0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutZip provides a mock function with given fields: mockeryArg0, mockeryArg1
func (mockerySelf *ZipStore) PutZip(mockeryArg0 coordinates.Module, mockeryArg1 repository.Blob) error {
	ret := mockerySelf.Called(mockeryArg0, mockeryArg1)

	var r0 error
	if rf, ok := ret.Get(0).(func(coordinates.Module, repository.Blob) error); ok {
		r0 = rf(mockeryArg0, mockeryArg1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
