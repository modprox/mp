package data

import (
	"database/sql"

	"github.com/pkg/errors"
)

const (
	insertModuleSQL = iota
	selectModuleIDSQL
	selectModulesBySource
	selectModuleIDScanSQL
	selectModulesByIDsSQL
	selectSourcesScanSQL
	insertHeartbeatSQL
	insertStartupConfigSQL
	selectStartupConfigsSQL
	selectHeartbeatsSQL
	deleteHeartbeatSQL
	deleteStartupConfigSQL
	deleteModuleByIDSQL
)

type statements map[int]*sql.Stmt

func load(kind string, db *sql.DB) (statements, error) {
	loaded := make(statements, len(mySQLTexts))

	stmtTexts := mySQLTexts
	if kind == "postgres" {
		return loaded, errors.New("postgres is not supports (issue #103)")
		// stmtTexts = postgreSQLTexts
	}

	for id, text := range stmtTexts {
		if text != "" { // avoid loading selectModulesByIDsSQL for mysql; must be generated
			stmt, err := db.Prepare(text)
			if err != nil {
				return nil, errors.Wrapf(err, "bad sql statement: %q", text)
			}
			loaded[id] = stmt
		}
	}

	return loaded, nil
}

var (
	mySQLTexts = map[int]string{
		insertModuleSQL:         `insert into modules(source, version) values (?, ?)`,
		selectModuleIDSQL:       `select id from modules where source=? and version=?`,
		selectModulesBySource:   `select id, source, version from modules where source=?`,
		selectModuleIDScanSQL:   `select id from modules order by id asc`,
		selectModulesByIDsSQL:   ``, // select id, source, version from modules where id in(?) order by id asc`,
		selectSourcesScanSQL:    `select id, source, version from modules`,
		insertHeartbeatSQL:      `insert into proxy_heartbeats (hostname, port, num_modules, num_versions) values (?, ?, ?, ?) on duplicate key update num_modules=?, num_versions=?, ts=current_timestamp;`,
		insertStartupConfigSQL:  `insert into proxy_configurations (hostname, port, storage, registry, transforms) values (?, ?, ?, ?, ?) on duplicate key update storage=?, registry=?, transforms=?`,
		selectStartupConfigsSQL: `select hostname, port, storage, registry, transforms from proxy_configurations`,
		selectHeartbeatsSQL:     `select hostname, port, num_modules, num_versions, unix_timestamp(ts) from proxy_heartbeats`,
		deleteHeartbeatSQL:      `delete from proxy_heartbeats where hostname=? and port=? limit 1`,
		deleteStartupConfigSQL:  `delete from proxy_configurations where hostname=? and port=? limit 1`,
		deleteModuleByIDSQL:     `delete from modules where id=?`,
	}

	/* issue #103 : put postgres on the shelf for now
	postgreSQLTexts = map[int]string{
		insertModuleSQL:         `insert into modules (source, version) values ($1, $2)`,
		selectModuleIDSQL:       `select id from modules where source=$1 and version=$2`,
		selectModulesBySource:   `select id, source, version from modules where source=$1`,
		selectModuleIDScanSQL:   `select id from modules order by id asc`,                                    // index scan module ids
		selectModulesByIDsSQL:   `select id, source, version from modules where id=any($1) order by id asc;`, // $1 is array
		selectSourcesScanSQL:    `select id, source, version from modules`,
		insertHeartbeatSQL:      `insert into proxy_heartbeats (hostname, port, num_modules, num_versions) values ($1, $2, $3, $4) on conflict (hostname, port) do update set num_modules=$5, num_versions=$6, ts=current_timestamp`,
		insertStartupConfigSQL:  `insert into proxy_configurations (hostname, port, storage, registry, transforms) values ($1, $2, $3, $4, $5) on conflict (hostname, port) do update set storage=$6, registry=$7, transforms=$8`,
		selectStartupConfigsSQL: `select hostname, port, storage, registry, transforms from proxy_configurations`,
		selectHeartbeatsSQL:     `select hostname, port, num_modules, num_versions, (cast( -extract(timezone from now()) + extract(epoch from ts) as integer)) from proxy_heartbeats`,
		deleteHeartbeatSQL:      `delete from proxy_heartbeats where hostname=$1 and port=$2`,
		deleteStartupConfigSQL:  `delete from proxy_configurations where hostname=$1 and port=$2`,
		// todo: deleteModuleByIDSQL
	}
	*/
)
