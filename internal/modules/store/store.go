package store

import "github.com/modprox/libmodprox/repository"

type Blob []byte

//go:generate mockery -interface=Store -package=storetest

type Store interface {
	List() ([]repository.ModInfo, error)
	Set(repository.ModInfo, Blob) error
	Get(repository.ModInfo) (Blob, error)
}
