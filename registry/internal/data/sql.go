package data

import (
	"database/sql"

	"github.com/pkg/errors"
)

const (
	insertModuleSQL = iota
	selectModuleIDSQL
	selectModuleIDScanSQL
	selectModulesByIDsSQL
	selectSourcesScanSQL
	insertStartupConfigSQL
	insertHeartbeatSQL
)

type statements map[int]*sql.Stmt

func load(kind string, db *sql.DB) (statements, error) {
	loaded := make(statements, len(mySQLText))

	stmtText := mySQLText
	if kind == "postgres" {
		stmtText = postgreSQLText
	}

	for id, text := range stmtText {
		stmt, err := db.Prepare(text)
		if err != nil {
			return nil, errors.Wrapf(err, "bad sql statement: %q", text)
		}
		loaded[id] = stmt
	}

	return loaded, nil
}

var (
	mySQLText = map[int]string{
		// todo: implement mysql queries
	}

	postgreSQLText = map[int]string{
		insertModuleSQL:       `insert into modules (source, version) values ($1, $2)`,
		selectModuleIDSQL:     `select id from modules where source=$1 and version=$2`,
		selectModuleIDScanSQL: `select id from modules order by id asc`,                                    // index scan module ids
		selectModulesByIDsSQL: `select id, source, version from modules where id=any($1) order by id asc;`, // $1 is array
		selectSourcesScanSQL:  `select id, source, version from modules`,
		insertHeartbeatSQL:    `insert into proxy_heartbeats (hostname, port, num_packages, num_modules) values ($1, $2, $3, $4) on conflict (hostname, port) do update set num_packages=$5, num_modules=$6, ts=current_timestamp`,
		// insertStartupConfigSQL`, // todo: write startup configs on startup
	}
)
