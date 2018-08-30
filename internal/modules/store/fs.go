package store

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/rboyer/safeio"

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
	if options.Directory == "" {
		panic("no directory set for store")
	}

	return &fsStore{
		options: options,
		log:     loggy.New("fs-store"),
	}
}

func (s *fsStore) Get(mod repository.ModInfo) (repository.Blob, error) {
	return nil, nil
}

func (s *fsStore) Put(mod repository.ModInfo, blob repository.Blob) error {
	s.log.Infof("will save %s do disk, %d bytes", mod, len(blob))
	exists, err := s.exists(mod)
	if err != nil {
		return err
	}

	if exists {
		s.log.Warnf("not saving %s because we already have it @ %s", mod, pathOf)
		return errors.Errorf("already have a copy of %s", mod)
	}

	// todo: these are not safe relative to each other; we need to create a tmp
	// todo: directory, write the files into that directory, atomic rename the
	// todo: directory, and then fsync the parent directory.

	if err := s.safeWriteZip(mod, blob); err != nil {
		s.log.Errorf("failed to write zip for %s, %v", mod, err)
		return err
	}

	if err := s.safeWriteModFile(mod, blob); err != nil {
		s.log.Errorf("failed to write go.mod for %s, %v", mod, err)
		return err
	}

	if err := s.safeWriteInfoFile(mod); err != nil {
		s.log.Errorf("failed to write .info for %s, %v", mod, err)
		return err
	}

	return nil
}

const (
	filePerm      = 0660
	directoryPerm = 0770
)

func (s *fsStore) safeWriteZip(mod repository.ModInfo, blob repository.Blob) error {
	modPath := s.fullPathOf(mod)
	s.log.Tracef("writing module zip into path: %s", modPath)

	// writing the zip always goes first, make sure the tree exists
	if err := os.MkdirAll(modPath, directoryPerm); err != nil {
		return err
	}

	zipFile := filepath.Join(modPath, zipName(mod))
	reader := bytes.NewReader(blob)
	_, err := safeio.WriteToFile(reader, zipFile, filePerm)
	return err
}

func (s *fsStore) safeWriteModFile(mod repository.ModInfo, blob repository.Blob) error {
	modPath := s.fullPathOf(mod)
	s.log.Tracef("writing module go.mod into path: %s", modPath)

	modFile := filepath.Join(modPath, modName(mod))
	content, exists, err := blob.ModFile()
	if err != nil {
		return err
	}

	if !exists {
		// fudge a blank go.mod file
		content = emptyModFile(mod)
	}

	reader := strings.NewReader(content)
	_, err = safeio.WriteToFile(reader, modFile, filePerm)
	return err
}

func (s *fsStore) safeWriteInfoFile(mod repository.ModInfo) error {
	modPath := s.fullPathOf(mod)
	s.log.Tracef("writing module .info into path: %s", modPath)

	infoFile := filepath.Join(modPath, infoName(mod))
	revInfo := newRevInfo(mod)
	content := revInfo.String()

	reader := strings.NewReader(content)
	_, err := safeio.WriteToFile(reader, infoFile, filePerm)
	return err
}

func zipName(mod repository.ModInfo) string {
	return mod.Version + ".zip"
}

func modName(mod repository.ModInfo) string {
	return mod.Version + ".mod"
}

func infoName(mod repository.ModInfo) string {
	return mod.Version + ".info"
}

func emptyModFile(mod repository.ModInfo) string {
	return fmt.Sprintf("module %s\n", mod.Source)
}

func newRevInfo(mod repository.ModInfo) repository.RevInfo {
	// todo: ... what goes in the revinfo?
	return repository.RevInfo{
		Version: mod.Version,
	}
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
