package store

import (
	"oss.indeed.com/go/modprox/pkg/coordinates"
	"oss.indeed.com/go/modprox/pkg/repository"
)

//go:generate mockery3 -interface=ZipStore -package=storetest

type ZipStore interface {
	PutZip(coordinates.Module, repository.Blob) error
	GetZip(coordinates.Module) (repository.Blob, error)
	DelZip(coordinates.Module) error
}
