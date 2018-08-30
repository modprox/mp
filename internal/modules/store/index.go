package store

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/modprox/libmodprox/loggy"
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
		log:     loggy.New("fs-index"),
	}
}

// fsIndex is an MVP implementation of Index which just
// retrieves all the information "live" from the actual
// filesystem store. This is very slow, but easy to implement.
type fsIndex struct {
	options IndexOptions
	log     loggy.Logger
}

func (i *fsIndex) List(module string) ([]string, error) {
	versionsDir := i.versionsPathOf(module)
	list, err := ioutil.ReadDir(versionsDir)
	if err != nil {
		i.log.Errorf("unable to list versions directory for %s, %v", module, err)
		return nil, err
	}

	zips := make([]string, 0, 10)
	for _, file := range list {
		if strings.HasSuffix(file.Name(), ".zip") {
			version := strings.TrimSuffix(file.Name(), ".zip")
			zips = append(zips, version)
		}
	}

	return zips, nil
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
