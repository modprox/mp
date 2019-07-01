package store

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"

	"github.com/shoenig/atomicfs"

	"go.gophers.dev/pkgs/loggy"

	"oss.indeed.com/go/modprox/pkg/coordinates"
	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/pkg/repository"
)

const (
	filePerm      = 0660
	directoryPerm = 0770
)

type fsStore struct {
	options Options
	emitter stats.Sender
	writer  atomicfs.FileWriter
	log     loggy.Logger
}

type Options struct {
	Directory    string
	TmpDirectory string
}

func NewStore(options Options, emitter stats.Sender) ZipStore {
	if options.Directory == "" {
		panic("no directory set for store")
	}

	writer := atomicfs.NewFileWriter(atomicfs.Options{
		TmpDirectory: options.TmpDirectory,
		Mode:         filePerm,
	})

	return &fsStore{
		options: options,
		emitter: emitter,
		writer:  writer,
		log:     loggy.New("fs-store"),
	}
}

func (s *fsStore) GetZip(mod coordinates.Module) (repository.Blob, error) {
	s.log.Tracef("retrieving module %s", mod)

	start := time.Now()
	blob, err := s.getZip(mod)
	if err != nil {
		s.emitter.Count("fsstore-getzip-failure", 1)
		return nil, err
	}

	s.emitter.GaugeMS("fsstore-getzip-elapsed-ms", start)
	return blob, nil
}

func (s *fsStore) getZip(mod coordinates.Module) (repository.Blob, error) {
	zipFile := filepath.Join(
		s.fullPathOf(mod),
		zipName(mod),
	)
	return ioutil.ReadFile(zipFile)
}

func (s *fsStore) DelZip(mod coordinates.Module) error {
	s.log.Tracef("removing module %s", mod)

	start := time.Now()
	err := s.removeZip(mod)
	if err != nil {
		s.emitter.Count("fsstore-rmzip-failure", 1)
		return err
	}

	s.emitter.GaugeMS("fsstore-rmzip-elapsed-ms", start)
	return nil
}

func (s *fsStore) removeZip(mod coordinates.Module) error {
	zipFile := filepath.Join(
		s.fullPathOf(mod),
		zipName(mod),
	)
	return os.Remove(zipFile)
}

func (s *fsStore) PutZip(mod coordinates.Module, blob repository.Blob) error {
	s.log.Infof("will save %s to disk, %d bytes", mod, len(blob))

	start := time.Now()
	if err := s.putZip(mod, blob); err != nil {
		s.emitter.Count("fsstore-putzip-failure", 1)
		return err
	}

	s.emitter.GaugeMS("fsstore-putzip-elapsed-ms", start)
	return nil
}

func (s *fsStore) putZip(mod coordinates.Module, blob repository.Blob) error {
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

func (s *fsStore) safeWriteZip(mod coordinates.Module, blob repository.Blob) error {
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

func zipName(mod coordinates.Module) string {
	return mod.Version + ".zip"
}

func (s *fsStore) exists(mod coordinates.Module) (bool, error) {
	modPath := s.fullPathOf(mod)
	_, err := os.Stat(modPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err != nil, err
}

func (s *fsStore) fullPathOf(mod coordinates.Module) string {
	return filepath.Join(
		s.options.Directory,
		pathOf(mod),
	)
}

func pathOf(mod coordinates.Module) string {
	return filepath.FromSlash(mod.Source) // eh windows?
}
