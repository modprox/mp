package store

import (
	"github.com/modprox/libmodprox/coordinates"
	"github.com/modprox/libmodprox/repository"
)

//go:generate mockery -interface=ZipStore -package=storetest

type ZipStore interface {
	PutZip(coordinates.Module, repository.Blob) error
	GetZip(coordinates.Module) (repository.Blob, error)
}
