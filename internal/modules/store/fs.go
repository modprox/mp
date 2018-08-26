package store

import "github.com/modprox/libmodprox/repository"

type fsStore struct {
	options Options
}

type Options struct {
	Directory string
}

func NewStore(options Options) Store {
	return &fsStore{
		options: options,
	}
}

func (s *fsStore) List() ([]repository.ModInfo, error) {
	return nil, nil
}

func (s *fsStore) Set(m repository.ModInfo, b Blob) error {
	return nil
}

func (s *fsStore) Get(m repository.ModInfo) (Blob, error) {
	return nil, nil
}
