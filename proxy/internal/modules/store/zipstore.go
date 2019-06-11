package store

import (
	"oss.indeed.com/go/modprox/pkg/coordinates"
	"oss.indeed.com/go/modprox/pkg/repository"
)

//go:generate minimock -g -i ZipStore -s _mock.go

type ZipStore interface {
	PutZip(coordinates.Module, repository.Blob) error
	GetZip(coordinates.Module) (repository.Blob, error)
	DelZip(coordinates.Module) error
}
