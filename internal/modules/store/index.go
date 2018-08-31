package store

import (
	"encoding/json"
	"io/ioutil"
	"os"
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

	revInfoPath := filepath.Join(
		i.versionsPathOf(mod.Source),
		mod.Version,
	) + ".info"
	i.log.Tracef("looking up revinfo at path %s", revInfoPath)

	f, err := os.Open(revInfoPath)
	if err != nil {
		i.log.Errorf("failed to open revinfo file at %s", revInfoPath)
		return revInfo, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&revInfo); err != nil {
		return revInfo, err
	}

	return revInfo, nil
}

func (i *fsIndex) Mod(mod repository.ModInfo) (string, error) {
	modFilePath := filepath.Join(
		i.versionsPathOf(mod.Source),
		mod.Version,
	) + ".mod"
	i.log.Tracef("looking up mod file at path %s", modFilePath)

	bs, err := ioutil.ReadFile(modFilePath)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
