package repositories

import (
	"database/sql"

	"github.com/pkg/errors"
)

const (
	selectAllRegistryInfos = iota
	insertRegistryInfo
)

type statements map[int]*sql.Stmt

func load(db *sql.DB) (statements, error) {
	loaded := make(statements, len(sqlText))

	for id, text := range sqlText {
		stmt, err := db.Prepare(text)
		if err != nil {
			return nil, errors.Wrapf(err, "bad sql statement: %q", text)
		}
		loaded[id] = stmt
	}

	return loaded, nil
}

var (
	sqlText = map[int]string{
		selectAllRegistryInfos: `select source, version from registry`,
		insertRegistryInfo:     `insert into registry (source, version) values (?, ?) on duplicate key update source=source`,
	}
)
