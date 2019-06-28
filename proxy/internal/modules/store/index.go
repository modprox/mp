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

	"oss.indeed.com/go/modprox/pkg/coordinates"
	"oss.indeed.com/go/modprox/pkg/loggy"
	"oss.indeed.com/go/modprox/pkg/repository"
)

// Ranges is an alias of coordinates.RangeIDs for brevity.
type Ranges = coordinates.RangeIDs

// Range is an alias of coordinates.RangeID for brevity.
type Range = coordinates.RangeID

//go:generate go run github.com/gojuno/minimock/cmd/minimock -g -i Index -s _mock.go

// The Index is used to provide:
//  - .mod file content
//  - .info file content
//  - boolean whether a module@version exists in the store
//  - list of versions of a given module that exist in the store
//  - list of version intervals for all modules in the store
//
// The real implementation is an index backed by boltdb, so
// we get better performance than keeping actual files on disk.
type Index interface {
	Versions(module string) ([]string, error)
	Info(coordinates.Module) (repository.RevInfo, error) // .info json
	Contains(coordinates.Module) (bool, int64, error)
	UpdateID(coordinates.SerialModule) error
	Mod(coordinates.Module) (string, error) // go.mod file
	Remove(coordinates.Module) error
	Put(ModuleAddition) error
	IDs() (Ranges, error)
	Summary() (int, int, error)
}

type ModuleAddition struct {
	Mod      coordinates.Module
	UniqueID int64
	ModFile  string
}

func (m ModuleAddition) String() string {
	return m.Mod.String()
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

func (i *boltIndex) Info(mod coordinates.Module) (repository.RevInfo, error) {
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

func (i *boltIndex) Mod(mod coordinates.Module) (string, error) {
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

func (i *boltIndex) Contains(mod coordinates.Module) (bool, int64, error) {
	key := mod.Bytes()
	var exists bool
	var id int64

	err := i.db.View(func(tx *bolt.Tx) error {
		idBkt := tx.Bucket(idBktLbl)
		bs := idBkt.Get(key)
		exists = bs != nil
		if exists {
			id = decodeID(bs)
		}
		return nil
	})

	return exists, id, err
}

func (i *boltIndex) Remove(mod coordinates.Module) error {
	key := mod.Bytes()
	return i.db.Update(func(tx *bolt.Tx) error {
		// need to remove the mod from each bucket

		// 1) remove from mods bucket
		if err := i.removeFromMods(key, tx); err != nil {
			return err
		}

		// 2) remove from info bucket
		if err := i.removeFromInfos(key, tx); err != nil {
			return err
		}

		// 3) remove from ids bucket
		if err := i.removeFromIDs(key, tx); err != nil {
			return err
		}

		return nil
	})
}

func (i *boltIndex) removeFromMods(key []byte, tx *bolt.Tx) error {
	modsBkt := tx.Bucket(modsBktLbl)
	return modsBkt.Delete(key)
}

func (i *boltIndex) removeFromInfos(key []byte, tx *bolt.Tx) error {
	infoBkt := tx.Bucket(infoBktLbl)
	return infoBkt.Delete(key)
}

func (i *boltIndex) removeFromIDs(key []byte, tx *bolt.Tx) error {
	idBkt := tx.Bucket(idBktLbl)
	return idBkt.Delete(key)
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
			encodedID := encodeID(add.UniqueID)
			idBkt := tx.Bucket(idBktLbl)
			if err := idBkt.Put(key, encodedID); err != nil {
				return err
			}
		}

		return nil
	})
}

func encodeID(id int64) []byte {
	var encodedID = make([]byte, 8) // 8 bytes in uint64
	binary.BigEndian.PutUint64(encodedID, uint64(id))
	return encodedID
}

func decodeID(bs []byte) int64 {
	uID := binary.BigEndian.Uint64(bs)
	return int64(uID)
}

func (i *boltIndex) UpdateID(mod coordinates.SerialModule) error {
	key := mod.Bytes()
	return i.db.Update(func(tx *bolt.Tx) error {
		encodedID := encodeID(mod.SerialID)
		idBkt := tx.Bucket(idBktLbl)
		// overwrite the ID for module
		if err := idBkt.Put(key, encodedID); err != nil {
			return err
		}
		return nil
	})
}

func newRevInfo(mod coordinates.Module) repository.RevInfo {
	// todo: ... what goes in the revinfo?
	return repository.RevInfo{
		Version: mod.Version,
	}
}

func (i *boltIndex) IDs() (Ranges, error) {
	var ids []int64 // values in the bucket

	err := i.db.View(func(tx *bolt.Tx) error {
		idBkt := tx.Bucket(idBktLbl)
		_ = idBkt.ForEach(func(_, v []byte) error {
			id := decodeID(v)
			ids = append(ids, id)
			return nil
		})

		return nil
	})

	return ranges(ids), err
}

func ranges(ids []int64) Ranges {
	var cuts Ranges

	sort.Slice(ids, func(x, y int) bool {
		return ids[x] < ids[y]
	})

	for {
		if len(ids) == 0 {
			return cuts
		}

		r, l := first(ids)
		cuts = append(cuts, r)
		ids = ids[l:]
	}
}

// just get the first sequence from ids
// this could be done without building the intermediate
// range, but meh (for now)
func first(ids []int64) (Range, int) {
	if len(ids) == 0 {
		return Range{0, 0}, 0
	}

	var seq []int64
	for i := 0; i < len(ids); i++ {
		if i == 0 {
			seq = append(seq, ids[i])
		} else if ids[i-1] == ids[i]-1 {
			seq = append(seq, ids[i])
		} else {
			break
		}
	}

	includes := Range{seq[0], seq[len(seq)-1]}
	return includes, len(seq)
}

func (i *boltIndex) Summary() (int, int, error) {
	m := make(map[string]int, 1024) // module => num versions

	if err := i.db.View(func(tx *bolt.Tx) error {
		idBkt := tx.Bucket(idBktLbl)
		return idBkt.ForEach(func(k, _ []byte) error {
			i := bytes.Index(k, []byte("@"))
			mod := string(k[:i])
			m[mod]++
			return nil
		})
	}); err != nil {
		return 0, 0, err
	}

	mods, versions := count(m)
	return mods, versions, nil
}

func count(m map[string]int) (int, int) {
	count := 0
	for _, v := range m {
		count += v
	}
	return len(m), count
}
