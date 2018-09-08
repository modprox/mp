package store

import (
	"github.com/modprox/libmodprox/repository"
)

//go:generate mockery -interface=ZipStore -package=storetest

type ZipStore interface {
	PutZip(repository.ModInfo, repository.Blob) error
	GetZip(repository.ModInfo) (repository.Blob, error)
}
