package store

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/shoenig/atomicfs"

	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/libmodprox/repository"
)

const (
	filePerm      = 0660
	directoryPerm = 0770
)

type fsStore struct {
	options Options
	writer  atomicfs.FileWriter
	log     loggy.Logger
}

type Options struct {
	Directory    string
	TmpDirectory string
}

func NewStore(options Options) ZipStore {
	if options.Directory == "" {
		panic("no directory set for store")
	}

	writer := atomicfs.NewFileWriter(atomicfs.Options{
		TmpDirectory: options.TmpDirectory,
		Mode:         filePerm,
	})

	return &fsStore{
		options: options,
		writer:  writer,
		log:     loggy.New("fs-store"),
	}
}

func (s *fsStore) GetZip(mod repository.ModInfo) (repository.Blob, error) {
	s.log.Tracef("retrieving module %s", mod)
	zipFile := filepath.Join(
		s.fullPathOf(mod),
		zipName(mod),
	)
	return ioutil.ReadFile(zipFile)
}

func (s *fsStore) PutZip(mod repository.ModInfo, blob repository.Blob) error {
	s.log.Infof("will save %s do disk, %d bytes", mod, len(blob))
	exists, err := s.exists(mod)
	if err != nil {
		return err
	}

	if exists {
		s.log.Warnf("not saving %s because we already have it @ %s", mod, pathOf)
		return errors.Errorf("already have a copy of %s", mod)
	}

	if err := s.safeWriteZip(mod, blob); err != nil {
		s.log.Errorf("failed to write zip for %s, %v", mod, err)
		return err
	}

	return nil
}

func (s *fsStore) safeWriteZip(mod repository.ModInfo, blob repository.Blob) error {
	modPath := s.fullPathOf(mod)
	s.log.Tracef("writing module zip into path: %s", modPath)

	// writing the zip always goes first, make sure the tree exists
	if err := os.MkdirAll(modPath, directoryPerm); err != nil {
		return err
	}

	zipFile := filepath.Join(modPath, zipName(mod))
	reader := bytes.NewReader(blob)
	return s.writer.Write(reader, zipFile)
}

func zipName(mod repository.ModInfo) string {
	return mod.Version + ".zip"
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
