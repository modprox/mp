package store

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/libmodprox/repository"
)

type fsStore struct {
	options Options
	log     loggy.Logger
}

type Options struct {
	Directory string
}

func NewStore(options Options) Store {
	return &fsStore{
		options: options,
		log:     loggy.New("fs-store"),
	}
}

func (s *fsStore) List() ([]repository.ModInfo, error) {
	return nil, nil
}

func (s *fsStore) Set(mod repository.ModInfo, b repository.Blob) error {
	s.log.Infof("will save %s do disk, %d bytes", mod, len(b))
	exists, err := s.exists(mod)
	if err != nil {
		return err
	}

	if exists {
		s.log.Warnf("not saving %s because we already have it @ %s", mod, pathOf)
		return errors.Errorf("already have a copy of %s", mod)
	}

	// todo, save it

	return nil
}

func (s *fsStore) Get(mod repository.ModInfo) (repository.Blob, error) {
	return nil, nil
}

func (s *fsStore) exists(mod repository.ModInfo) (bool, error) {
	modPath := s.fullPathOf(mod)
	_, err := os.Stat(modPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err != nil, err
}

func (s *fsStore) fullPathOf(mod repository.ModInfo) string {
	return filepath.Join(
		s.options.Directory,
		pathOf(mod),
	)
}

func pathOf(mod repository.ModInfo) string {
	return filepath.FromSlash(mod.Source) // eh windows?
}

/*
$ pwd /home/hoenig/Documents/go/pkg/mod/cache/download/github.com/pkg/errors/@v
$ cat -n list v0.8.0.info v0.8.0.mod v0.8.0.ziphash
     1	v0.8.0
     2	{"Version":"v0.8.0","Time":"2016-09-29T01:48:01Z"}module github.com/pkg/errors
     3	h1:WdK/asTD0HN+q6hsWO3/vpuAkAr+tw6aNJNDFFf0+qw=
*/
