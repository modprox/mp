package store

import (
	"github.com/modprox/libmodprox/repository"
)

//go:generate mockery -interface=Store -package=storetest

type Store interface {
	List() ([]repository.ModInfo, error)
	Set(repository.ModInfo, repository.Blob) error
	Get(repository.ModInfo) (repository.Blob, error)
}
