package store

import (
	"github.com/modprox/mp/pkg/coordinates"
	"github.com/modprox/mp/pkg/repository"
)

//go:generate mockery -interface=ZipStore -package=storetest

type ZipStore interface {
	PutZip(coordinates.Module, repository.Blob) error
	GetZip(coordinates.Module) (repository.Blob, error)
}
