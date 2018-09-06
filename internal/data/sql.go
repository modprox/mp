package data

import (
	"database/sql"

	"github.com/pkg/errors"
)

const (
	insertSourceSQL = iota
	insertTagSQL
	selectSourcesScanSQL
	insertRedirectSQL
	selectRedirectsScanSQL
	insertStartupConfigSQL
	insertHeartbeatSQL
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
		insertSourceSQL:        `insert into sources (source) values (?) on duplicate key update id=last_insert_id(id), source=source`,
		insertTagSQL:           `insert into tags (tag, source_id) values (?, ?) on duplicate key update id=last_insert_id(id), tag=tag`,
		selectSourcesScanSQL:   `select sources.id, sources.source, unix_timestamp(sources.created), tags.id, tags.tag, unix_timestamp(tags.created), tags.source_id from sources inner join (tags) on (tags.source_id=sources.id) order by sources.source asc, tags.tag desc`,
		insertRedirectSQL:      `insert into redirects (original, substitution) values (?, ?) on duplicate key update original=original`,
		selectRedirectsScanSQL: `select original, substitution from redirects order by original asc`,
		insertStartupConfigSQL: `insert into proxy_configurations (hostname, port, transforms) values (?, ?, ?) on duplicate key update id=last_insert_id(id), transforms=?, ts=current_timestamp`,
		insertHeartbeatSQL:     `insert into proxy_heartbeats (hostname, port, num_packages, num_modules) values (?, ?, ?, ?) on duplicate key update id=last_insert_id(id), num_packages=?, num_modules=?, ts=current_timestamp`,
	}
)
