package store

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"

	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/libmodprox/repository"
)

// the index is used to return content/info:
// - .mod files
// - .info files
// - boolean whether a module@version is in the store
// - list of versions in the store for a module
//
// in the future, we will need to
// - know the uniqueID of the module, assigned by the registry
//   for the use case of giving the registry an optimized gap-list
//   describing which modules we already have, so the registry
//   can respond with a minimized list of modules the proxy needs
//   to download

type Index interface {
	Versions(module string) ([]string, error)
	Info(repository.ModInfo) (repository.RevInfo, error)
	Mod(repository.ModInfo) (string, error) // go.mod
	Contains(repository.ModInfo) (bool, error)
	Put(ModuleAddition) error
}

type ModuleAddition struct {
	Mod      repository.ModInfo
	UniqueID uint64
	ModFile  string
}

type IndexOptions struct {
	Directory   string
	OpenTimeout time.Duration
}

func NewIndex(options IndexOptions) (Index, error) {
	log := loggy.New("bolt-index")

	if options.Directory == "" {
		return nil, errors.New("no directory set for index")
	}

	openTimeout := options.OpenTimeout
	if openTimeout <= 0 {
		openTimeout = 10 * time.Second
	}

	if err := setupDirs(options.Directory); err != nil {
		return nil, errors.Wrap(err, "unable to make directories for modprox.db")
	}

	dbPath := filepath.Join(options.Directory, "modprox.db")
	db, err := bolt.Open(dbPath, 0660, &bolt.Options{
		Timeout: openTimeout,
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to open modprox.db")
	}

	if err := initDB(db); err != nil {
		return nil, errors.Wrap(err, "unable to initialize modprox.db")
	}

	return &boltIndex{
		options: options,
		db:      db,
		log:     log,
	}, nil
}

var (
	modsBktLbl = []byte("mods")
	infoBktLbl = []byte("info")
	idBktLbl   = []byte("ids")
)

func setupDirs(indexPath string) error {
	return os.MkdirAll(indexPath, 0770)
}

func initDB(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(modsBktLbl)); err != nil {
			return err
		}

		if _, err := tx.CreateBucketIfNotExists([]byte(infoBktLbl)); err != nil {
			return err
		}

		if _, err := tx.CreateBucketIfNotExists([]byte(idBktLbl)); err != nil {
			return err
		}
		return nil
	})
}

type boltIndex struct {
	options IndexOptions
	db      *bolt.DB
	log     loggy.Logger
}

func (i *boltIndex) Versions(module string) ([]string, error) {
	// produces an ordered list of version strings

	prefix := []byte(module)
	var versions []string

	err := i.db.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket(idBktLbl).Cursor()
		for key, _ := cursor.Seek(prefix); key != nil && bytes.HasPrefix(key, prefix); key, _ = cursor.Next() {
			version := versionOf(key)
			versions = append(versions, version)
		}
		return nil
	})

	// todo: sort versions using common lib
	sort.Strings(versions) // incorrect

	return versions, err
}

func versionOf(key []byte) string {
	s := string(key)
	vIdx := strings.Index(s, "@")
	return s[vIdx+1:]
}

func (i *boltIndex) Info(mod repository.ModInfo) (repository.RevInfo, error) {
	key := mod.Bytes()
	var revInfo repository.RevInfo
	var content []byte

	if err := i.db.View(func(tx *bolt.Tx) error {
		infoBkt := tx.Bucket(infoBktLbl)
		bs := infoBkt.Get(key) // must copy inside tx
		content = make([]byte, len(bs))
		copy(content, bs)
		if bs == nil {
			return errors.New("module not in index")
		}
		return nil
	}); err != nil {
		return revInfo, err
	}

	err := json.Unmarshal(content, &revInfo)
	return revInfo, err
}

func (i *boltIndex) Mod(mod repository.ModInfo) (string, error) {
	key := mod.Bytes()
	var content string

	err := i.db.View(func(tx *bolt.Tx) error {
		modBkt := tx.Bucket(modsBktLbl)
		bs := modBkt.Get(key)
		content = string(bs)
		if bs == nil {
			return errors.New("module not in index")
		}
		return nil
	})

	return content, err
}

func (i *boltIndex) Contains(mod repository.ModInfo) (bool, error) {
	key := mod.Bytes()
	var exists bool

	err := i.db.View(func(tx *bolt.Tx) error {
		idBkt := tx.Bucket(idBktLbl)
		bs := idBkt.Get(key)
		exists = bs != nil
		return nil
	})

	return exists, err
}

func (i *boltIndex) Put(add ModuleAddition) error {
	key := add.Mod.Bytes()

	// update the three buckets with the new information
	return i.db.Update(func(tx *bolt.Tx) error {
		// insert the .mod file
		{
			modFile := []byte(add.ModFile)
			modsBkt := tx.Bucket(modsBktLbl)
			if err := modsBkt.Put(key, modFile); err != nil {
				return err
			}
		}

		// insert the .info file
		{
			infoFile := newRevInfo(add.Mod).Bytes()
			infoBkt := tx.Bucket(infoBktLbl)
			if err := infoBkt.Put(key, infoFile); err != nil {
				return err
			}
		}

		// insert the uniqueID
		{
			var encodedID = make([]byte, 8) // 8 bytes in uint64
			binary.BigEndian.PutUint64(encodedID, add.UniqueID)
			idBkt := tx.Bucket(idBktLbl)
			if err := idBkt.Put(key, encodedID); err != nil {
				return err
			}
		}

		return nil
	})
}

func newRevInfo(mod repository.ModInfo) repository.RevInfo {
	// todo: ... what goes in the revinfo?
	return repository.RevInfo{
		Version: mod.Version,
	}
}
