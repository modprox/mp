package repositories

import (
	"database/sql"

	"github.com/pkg/errors"
)

const (
	insertSourceSQL = iota
	insertTagSQL
	selectSourcesScanSQL
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
		insertSourceSQL:      `insert into sources (source) values (?) on duplicate key update source=source`,
		insertTagSQL:         `insert into tags (tag, source_id) values (?, ?) on duplicate key update tag=tag`,
		selectSourcesScanSQL: `select sources.id, sources.source, unix_timestamp(sources.created), tags.id, tags.tag, unix_timestamp(tags.created), tags.source_id from sources inner join (tags) on (tags.source_id=sources.id) order by sources.source asc, tags.tag desc`,
	}
)
