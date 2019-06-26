package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"oss.indeed.com/go/modprox/pkg/coordinates"
	"oss.indeed.com/go/modprox/pkg/loggy"
	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/pkg/repository"
)

type mysqlStore struct {
	emitter    stats.Sender
	db         *sql.DB
	statements statements
	log        loggy.Logger
}

var _ ZipStore = (*mysqlStore)(nil)
var _ Index = (*mysqlStore)(nil)

const dbTimeout = 10 * time.Second

func (m *mysqlStore) PutZip(mod coordinates.Module, blob repository.Blob) error {
	start := time.Now()
	err := m.insertModulesZip(mod, blob)
	if err != nil {
		m.emitter.Count("db-put-zip-failure", 1)
		return err
	}

	m.emitter.GaugeMS("db-put-zip-elapsed-ms", start)
	return nil
}

func (m *mysqlStore) GetZip(mod coordinates.Module) (repository.Blob, error) {
	m.log.Tracef("retrieving module %s", mod)

	start := time.Now()
	blob, err := m.getModuleZip(mod)
	if err != nil {
		return nil, err
	}

	m.emitter.GaugeMS("db-getzip-elapsed-ms", start)
	return blob, nil
}

func (m *mysqlStore) DelZip(mod coordinates.Module) error {
	m.log.Tracef("removing module %+v", mod)

	start := time.Now()
	err := m.removeModuleZip(mod)
	if err != nil {
		m.emitter.Count("db-rmzip-failure", 1)
		return err
	}

	m.emitter.GaugeMS("db-rmzip-elapsed-ms", start)
	return nil
}

func (m *mysqlStore) Versions(mod string) ([]string, error) {
	m.log.Tracef("retrieving versions for module %s", mod)

	start := time.Now()
	versions, err := m.getModuleVersions(mod)
	if err != nil {
		return nil, err
	}

	m.emitter.GaugeMS("db-get-versions-elapsed-ms", start)
	return versions, nil
}

func (m *mysqlStore) Info(mod coordinates.Module) (repository.RevInfo, error) {
	m.log.Tracef("retrieving revinfo for module %s", mod)

	start := time.Now()
	revInfo, err := m.getVersionInfo(mod)
	if err != nil {
		return repository.RevInfo{}, err
	}

	m.emitter.GaugeMS("db-get-revinfo-elapsed-ms", start)
	return revInfo, nil
}

func (m *mysqlStore) Contains(mod coordinates.Module) (bool, int64, error) {
	m.log.Tracef("retrieving registry ID for module %s", mod)

	start := time.Now()
	exists, id, err := m.getRegistryID(mod)
	if err == nil {
		m.emitter.GaugeMS("db-get-regid-elapsed-ms", start)
	}

	return exists, id, err
}

func (m *mysqlStore) UpdateID(mod coordinates.SerialModule) error {
	m.log.Tracef("updating registry id for module %s", mod)

	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	_, err := m.statements[updateRegistryIDSQL].ExecContext(ctx, mod.SerialID, mod.Source, mod.Version)
	if err != nil {
		m.emitter.Count("db-update-regid-failure", 1)
	} else {
		m.emitter.GaugeMS("db-update-regid-elapsed-ms", start)
	}

	return err
}

func (m *mysqlStore) Mod(mod coordinates.Module) (string, error) {
	m.log.Tracef("retrieving mod for module %s", mod)

	start := time.Now()
	revInfo, err := m.getGoMod(mod)
	if err != nil {
		return "", err
	}

	m.emitter.GaugeMS("db-get-mod-elapsed-ms", start)
	return revInfo, nil
}

func (m *mysqlStore) Remove(mod coordinates.Module) error {
	m.log.Tracef("deleting module %s", mod)

	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	_, err := m.statements[deleteModuleSQL].ExecContext(ctx, mod.Source, mod.Version)
	if err != nil {
		m.emitter.Count("db-delete-mod-failure", 1)
	} else {
		m.emitter.GaugeMS("db-delete-mod-elapsed-ms", start)
	}

	return err
}

func (m *mysqlStore) Put(add ModuleAddition) error {
	m.log.Tracef("adding module %s", add)

	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	_, err := m.statements[insertModuleSQL].
		ExecContext(ctx, add.Mod.Source, add.Mod.Version, []byte(add.ModFile), newRevInfo(add.Mod).Bytes(), add.UniqueID)
	if err != nil {
		m.emitter.Count("db-insert-mod-failure", 1)
	} else {
		m.emitter.GaugeMS("db-insert-mod-elapsed-ms", start)
	}

	return err
}

func (m *mysqlStore) IDs() (Ranges, error) {
	m.log.Tracef("retrieving all ids")

	start := time.Now()
	ids, err := m.ids()
	if err != nil {
		return nil, err
	}

	m.emitter.GaugeMS("db-get-ids-elapsed-ms", start)
	return ranges(ids), nil
}

func (m *mysqlStore) Summary() (int, int, error) {
	m.log.Tracef("retrieving all sources and computing summary")
	return m.countSourcesAndVersions()
}

const (
	insertModuleZipSQL = iota
	selectModuleZipSQL
	zipExistsSQL
	deleteModuleZipSQL
	insertModuleSQL
	selectRegistryIDSQL
	selectAllRegistryIDsSQL
	countVersionsSQL
	selectModuleVersionInfoSQL
	selectGoModFileSQL
	selectModuleVersionsSQL
	updateRegistryIDSQL
	deleteModuleSQL
)

type statements map[int]*sql.Stmt

func load(kind string, db *sql.DB) (statements, error) {
	loaded := make(statements, len(mySQLTexts))

	stmtTexts := mySQLTexts
	if kind == "postgres" {
		return loaded, errors.New("postgres is not supported (issue #103)")
		// stmtTexts = postgreSQLTexts
	}

	for id, text := range stmtTexts {
		stmt, err := db.Prepare(text)
		if err != nil {
			return nil, errors.Wrapf(err, "bad sql statement: %q", text)
		}
		loaded[id] = stmt
	}

	return loaded, nil
}

var (
	mySQLTexts = map[int]string{
		// modules
		insertModuleZipSQL: `insert into proxy_module_zips(path, zip) values (?, ?)`,
		selectModuleZipSQL: `select zip from proxy_module_zips where path=?`,
		zipExistsSQL:       `select count(id) from proxy_module_zips where path=?`,
		deleteModuleZipSQL: `delete from proxy_module_zips where path=?`,
		// index
		insertModuleSQL:            `insert into proxy_modules_index(source, version, go_mod_file, version_info, registry_mod_id) values (?, ?, ?, ?, ?)`,
		selectRegistryIDSQL:        `select registry_mod_id from proxy_modules_index where source=? and version=?`,
		selectAllRegistryIDsSQL:    `select registry_mod_id from proxy_modules_index`,
		countVersionsSQL:           `select count(version) from proxy_modules_index group by source`,
		selectModuleVersionInfoSQL: `select version_info from proxy_modules_index where source=? and version=?`,
		selectGoModFileSQL:         `select go_mod_file from proxy_modules_index where source=? and version=?`,
		selectModuleVersionsSQL:    `select version from proxy_modules_index where source=?`,
		updateRegistryIDSQL:        `update proxy_modules_index set registry_mod_id=? where source=? and version=?`,
		deleteModuleSQL:            `delete from proxy_modules_index where source=? and version=?`,
	}
)

func (m *mysqlStore) insertModulesZip(mod coordinates.Module, blob repository.Blob) error {
	exists, err := m.zipExists(mod)
	if err != nil {
		return err
	}

	path := pathOf(mod)
	if exists {
		m.log.Warnf("not saving %s because we already have it @ %s", mod, path)
		return errors.Errorf("already have a copy of %s", mod)
	}
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	if _, err = m.statements[insertModuleZipSQL].ExecContext(ctx, path, []byte(blob)); err != nil {
		m.emitter.Count("db-insertmodule-failure", 1)
		m.log.Errorf("failed to write zip for %s, %+v", mod, err)
	}
	return err
}

func (m *mysqlStore) zipExists(mod coordinates.Module) (bool, error) {
	path := pathOf(mod)
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	rows, err := m.statements[zipExistsSQL].QueryContext(ctx, path)
	if err != nil {
		m.emitter.Count("db-module-exists-failure", 1)
		return false, err
	}
	defer rows.Close()

	if !rows.Next() {
		m.emitter.Count("db-module-exists-failure", 1)
		return false, fmt.Errorf("expected exactly one row for sql: %+v", m.statements[zipExistsSQL])
	}

	var count int64
	err = rows.Scan(&count)
	if err != nil {
		m.emitter.Count("db-module-exists-failure", 1)
		return false, errors.Wrapf(err, "failed to read row for sql: %+v", m.statements[zipExistsSQL])
	}

	m.emitter.GaugeMS("db-module-exists-elapsed-ms", start)
	return count > 0, nil
}

func (m *mysqlStore) getModuleZip(mod coordinates.Module) (repository.Blob, error) {
	path := pathOf(mod)
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	rows, err := m.statements[selectModuleZipSQL].QueryContext(ctx, path)
	if err != nil {
		m.emitter.Count("db-select-module-failure", 1)
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		m.emitter.Count("db-select-module-failure", 1)
		return nil, fmt.Errorf("expected exactly one row for sql: %+v", m.statements[selectModuleZipSQL])
	}
	var contents []byte
	err = rows.Scan(&contents)
	if err != nil {
		m.emitter.Count("db-select-module-failure", 1)
		return nil, errors.Wrapf(err, "failed to read row for sql: %+v", m.statements[selectModuleZipSQL])
	}

	return repository.Blob(contents), nil
}

func (m *mysqlStore) removeModuleZip(mod coordinates.Module) error {
	path := pathOf(mod)
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	res, err := m.statements[deleteModuleZipSQL].ExecContext(ctx, path)
	if err != nil {
		m.emitter.Count("db-delete-module-failure", 1)
		return errors.Wrapf(err, "failed to delete zip for %s", mod)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		m.emitter.Count("db-delete-module-failure", 1)
		return errors.Wrapf(err, "failed to test rows affected for %s", mod)
	}
	if rowsAffected != 1 {
		return errors.Errorf("expected exactly 1 row to be deleted for %s but got %+v", mod, rowsAffected)
	}
	return nil
}

func (m *mysqlStore) getModuleVersions(source string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	rows, err := m.statements[selectModuleVersionsSQL].QueryContext(ctx, source)
	if err != nil {
		m.emitter.Count("db-select-module-versions-failure", 1)
		return nil, errors.Wrapf(err, "failed to query versions for %s", source)
	}
	defer rows.Close()

	versions := make([]string, 0, 10)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			m.emitter.Count("db-select-module-versions-failure", 1)
			return nil, errors.Wrapf(err, "failed to scan row for sql: %+v", m.statements[selectModuleVersionsSQL])
		}
		versions = append(versions, version)
	}

	if err := rows.Err(); err != nil {
		m.emitter.Count("db-select-module-versions-failure", 1)
		return nil, errors.Wrapf(err, "got error from rows for %s", source)
	}

	return versions, nil
}

func (m *mysqlStore) getVersionInfo(mod coordinates.Module) (repository.RevInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	rows, err := m.statements[selectModuleVersionInfoSQL].QueryContext(ctx, mod.Source, mod.Version)
	var revInfo repository.RevInfo
	if err != nil {
		m.emitter.Count("db-select-revinfo-failure", 1)
		return revInfo, errors.Wrapf(err, "failed to query revinfo for %+v", mod)
	}
	defer rows.Close()

	if !rows.Next() {
		return revInfo, errors.New("module not in index")
	}

	var contents []byte
	err = rows.Scan(&contents)
	if err != nil {
		m.emitter.Count("db-select-revinfo-failure", 1)
		return revInfo, errors.Wrapf(err, "failed to read row for sql: %+v", m.statements[selectModuleVersionInfoSQL])
	}

	err = json.Unmarshal(contents, &revInfo)
	return revInfo, err
}

func (m *mysqlStore) getRegistryID(mod coordinates.Module) (bool, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	rows, err := m.statements[selectRegistryIDSQL].QueryContext(ctx, mod.Source, mod.Version)
	if err != nil {
		m.emitter.Count("db-select-regid-failure", 1)
		return false, 0, errors.Wrapf(err, "failed to query regid for %+v", mod)
	}
	defer rows.Close()

	if !rows.Next() {
		return false, 0, nil
	}

	var regid int64
	err = rows.Scan(&regid)
	if err != nil {
		m.emitter.Count("db-select-regid-failure", 1)
		return false, 0, errors.Wrapf(err, "failed to read row for sql: %+v", m.statements[selectRegistryIDSQL])
	}

	return true, regid, nil
}

func (m *mysqlStore) getGoMod(mod coordinates.Module) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	rows, err := m.statements[selectGoModFileSQL].QueryContext(ctx, mod.Source, mod.Version)
	if err != nil {
		m.emitter.Count("db-select-gomod-failure", 1)
		return "", errors.Wrapf(err, "failed to query go.mod for %+v", mod)
	}
	defer rows.Close()

	if !rows.Next() {
		return "", errors.New("module not in index")
	}

	var contents []byte
	err = rows.Scan(&contents)
	if err != nil {
		m.emitter.Count("db-select-gomod-failure", 1)
		return "", errors.Wrapf(err, "failed to read row for sql: %+v", m.statements[selectGoModFileSQL])
	}

	return string(contents), err
}

func (m *mysqlStore) ids() ([]int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	rows, err := m.statements[selectAllRegistryIDsSQL].QueryContext(ctx)
	if err != nil {
		m.emitter.Count("db-select-ids-failure", 1)
		return nil, errors.Wrapf(err, "failed to query ids")
	}
	defer rows.Close()

	ids := make([]int64, 0, 10)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			m.emitter.Count("db-select-ids-failure", 1)
			return nil, errors.Wrapf(err, "failed to scan row for sql: %+v", m.statements[selectAllRegistryIDsSQL])
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		m.emitter.Count("db-select-ids-failure", 1)
		return nil, errors.Wrapf(err, "got error from rows")
	}

	return ids, nil
}

func (m *mysqlStore) countSourcesAndVersions() (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	rows, err := m.statements[countVersionsSQL].QueryContext(ctx)
	if err != nil {
		m.emitter.Count("db-count-versions-failure", 1)
		return 0, 0, errors.Wrapf(err, "failed to query sources")
	}
	defer rows.Close()

	totalVersions := 0
	totalSources := 0
	for rows.Next() {
		var count int
		if err := rows.Scan(&count); err != nil {
			m.emitter.Count("db-count-versions-failure", 1)
			return 0, 0, errors.Wrapf(err, "failed to scan row for sql: %+v", m.statements[countVersionsSQL])
		}
		totalSources += 1
		totalVersions += count
	}

	if err := rows.Err(); err != nil {
		m.emitter.Count("db-count-versions-failure", 1)
		return 0, 0, errors.Wrapf(err, "got error from rows")
	}

	return totalSources, totalVersions, nil
}
