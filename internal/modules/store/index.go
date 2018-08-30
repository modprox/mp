package store

import (
	"path/filepath"

	"github.com/modprox/libmodprox/repository"
)

type Index interface {
	List(module string) ([]string, error)
	Info(repository.ModInfo) (repository.RevInfo, error)
	Mod(repository.ModInfo) (string, error) // go.mod
}

type IndexOptions struct {
	Directory string
}

func NewIndex(options IndexOptions) Index {
	if options.Directory == "" {
		panic("no directory set for index")
	}

	return &fsIndex{
		options: options,
	}
}

// fsIndex is an MVP implementation of Index which just
// retrieves all the information "live" from the actual
// filesystem store. This is very slow, but easy to implement.
type fsIndex struct {
	options IndexOptions
}

func (i *fsIndex) List(module string) ([]string, error) {
	// versionsDir := i.versionsPathOf(module)
	// filepath.
	return nil, nil
}

func (i *fsIndex) versionsPathOf(module string) string {
	return filepath.Join(
		i.options.Directory,
		filepath.FromSlash(module),
	)
}

func (i *fsIndex) Info(mod repository.ModInfo) (repository.RevInfo, error) {
	var revInfo repository.RevInfo
	return revInfo, nil
}

func (i *fsIndex) Mod(mod repository.ModInfo) (string, error) {
	return "", nil
}
