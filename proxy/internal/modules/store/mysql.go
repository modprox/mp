package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"time"

	"github.com/pkg/errors"

	"oss.indeed.com/go/modprox/pkg/coordinates"
	"oss.indeed.com/go/modprox/pkg/loggy"
	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/pkg/repository"
)

type mysqlStore struct {
	statements statements
	db         *sql.DB
	emitter    stats.Sender
	log        loggy.Logger
}

var _ ZipStore = (*mysqlStore)(nil)
var _ Index = (*mysqlStore)(nil)

const dbTimeout = 10 * time.Second

// PutZip implements ZipStore.PutZip
func (m *mysqlStore) PutZip(mod coordinates.Module, blob repository.Blob) error {
	m.log.Tracef("put module zip %s", mod)
	start := time.Now()

	if err := m.insertModulesZip(mod, blob); err != nil {
		m.emitter.Count("db-put-zip-failure", 1)
		return err
	}

	m.emitter.GaugeMS("db-put-zip-elapsed-ms", start)
	return nil
}

// GetZip implements ZipStore.GetZip
func (m *mysqlStore) GetZip(mod coordinates.Module) (repository.Blob, error) {
	m.log.Tracef("get module zip %s", mod)
	start := time.Now()

	blob, err := m.getModuleZip(mod)
	if err != nil {
		m.emitter.Count("db-get-zip-failure", 1)
		return nil, err
	}

	m.emitter.GaugeMS("db-get-zip-elapsed-ms", start)
	return blob, nil
}

// DelZip implements ZipStore.DelZip
func (m *mysqlStore) DelZip(mod coordinates.Module) error {
	m.log.Tracef("del module zip %s", mod)
	start := time.Now()

	err := m.removeModuleZip(mod)
	if err != nil {
		m.emitter.Count("db-del-zip-failure", 1)
		return err
	}

	m.emitter.GaugeMS("db-del-zip-elapsed-ms", start)
	return nil
}

// Versions implements Index.Versions
func (m *mysqlStore) Versions(mod string) ([]string, error) {
	m.log.Tracef("get versions for module %s", mod)
	start := time.Now()

	versions, err := m.getModuleVersions(mod)
	if err != nil {
		m.emitter.Count("db-versions-failure", 1)
		return nil, err
	}

	m.emitter.GaugeMS("db-get-versions-elapsed-ms", start)
	return versions, nil
}

// Info implements Index.Info
func (m *mysqlStore) Info(mod coordinates.Module) (repository.RevInfo, error) {
	m.log.Tracef("get .info for module %s", mod)
	start := time.Now()

	modInfo, err := m.getVersionInfo(mod)
	if err != nil {
		m.emitter.Count("db-info-failure", 1)
		return repository.RevInfo{}, err
	}

	m.emitter.GaugeMS("db-info-elapsed-ms", start)
	return modInfo, nil
}

// Contains implements Index.Contains
func (m *mysqlStore) Contains(mod coordinates.Module) (bool, int64, error) {
	m.log.Tracef("get registry ID for module %s", mod)
	start := time.Now()

	exists, id, err := m.getRegistryID(mod)
	if err != nil {
		m.emitter.Count("db-contains-failure", 1)
		return false, 0, err
	}

	m.emitter.GaugeMS("db-get-registry-id-elapsed-ms", start)
	return exists, id, nil
}

// UpdateID implements Index.UpdateID
func (m *mysqlStore) UpdateID(mod coordinates.SerialModule) error {
	m.log.Tracef("updating registry id for module %s", mod)
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	_, err := m.statements[updateRegistryIDSQL].ExecContext(
		ctx,
		mod.SerialID,
		mod.Source,
		mod.Version,
	)
	if err != nil {
		m.emitter.Count("db-update-registry-id-failure", 1)
	} else {
		m.emitter.GaugeMS("db-update-registry-id-elapsed-ms", start)
	}

	return err
}

// Mod implements Index.Mod
func (m *mysqlStore) Mod(mod coordinates.Module) (string, error) {
	m.log.Tracef("retrieving mod for module %s", mod)
	start := time.Now()

	goMod, err := m.getGoMod(mod)
	if err != nil {
		m.emitter.Count("db-get-go-mod-failure", 1)
		return "", err
	}

	m.emitter.GaugeMS("db-get-mod-elapsed-ms", start)
	return goMod, nil
}

// Remove implements Index.Remove
func (m *mysqlStore) Remove(mod coordinates.Module) error {
	m.log.Tracef("deleting module %s", mod)
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	_, err := m.statements[deleteModuleSQL].ExecContext(
		ctx,
		mod.Source,
		mod.Version,
	)
	if err != nil {
		m.emitter.Count("db-delete-mod-failure", 1)
	} else {
		m.emitter.GaugeMS("db-delete-mod-elapsed-ms", start)
	}

	return err
}

// Put implements Index.Put
func (m *mysqlStore) Put(add ModuleAddition) error {
	m.log.Tracef("adding module %s", add)
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	_, err := m.statements[insertModuleSQL].ExecContext(
		ctx,
		add.Mod.Source,
		add.Mod.Version,
		[]byte(add.ModFile),
		newRevInfo(add.Mod).Bytes(),
		add.UniqueID,
	)
	if err != nil {
		m.emitter.Count("db-insert-mod-failure", 1)
	} else {
		m.emitter.GaugeMS("db-insert-mod-elapsed-ms", start)
	}

	return err
}

// IDs implements Index.IDs
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

// Summary implements Index.Summary
func (m *mysqlStore) Summary() (int, int, error) {
	m.log.Tracef("retrieving all sources and computing summary")
	return m.countSourcesAndVersions()
}

const (
	insertModuleZipSQL = iota
	selectModuleZipSQL
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

func load(db *sql.DB) (statements, error) {
	loaded := make(statements, len(mySQLTexts))

	stmtTexts := mySQLTexts
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
		// Do not reference registry tables since they are likely in a
		// separate database altogether.

		// Transactions are not used between the proxy_module_zips and
		// proxy_modules_index tables, because the proxy itself is designed to
		// keep these two discrete data-stores eventually consistent. This is
		// a historical feature due to the boltDB + filesystem implementation
		// that came first.

		// Table proxy_module_zips used to implement ZipStore.
		insertModuleZipSQL: `insert into proxy_module_zips(s_at_v, zip) values (?, ?)`,
		selectModuleZipSQL: `select zip from proxy_module_zips where s_at_v=?`,
		deleteModuleZipSQL: `delete from proxy_module_zips where s_at_v=?`,

		// Table proxy_modules_index used to implement Index.
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

// for ZipStore.PutZip
func (m *mysqlStore) insertModulesZip(mod coordinates.Module, blob repository.Blob) error {
	sAtV := mod.AtVersion()

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if _, err := m.statements[insertModuleZipSQL].ExecContext(ctx, sAtV, []byte(blob)); err != nil {
		m.emitter.Count("db-insert-module-failure", 1)
		m.log.Errorf("failed to write zip for %s, %+v", mod, err)
		return err
	}
	return nil
}

// for ZipStore.GetZip
func (m *mysqlStore) getModuleZip(mod coordinates.Module) (repository.Blob, error) {
	sAtV := mod.AtVersion()

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	row := m.statements[selectModuleZipSQL].QueryRowContext(ctx, sAtV)

	var contents []byte
	if err := row.Scan(&contents); err != nil {
		m.emitter.Count("db-select-module-failure", 1)
		return nil, errors.Wrapf(err, "failed to read row for sql: %+v", m.statements[selectModuleZipSQL])
	}

	return repository.Blob(contents), nil
}

// for ZipStore.DelZip
func (m *mysqlStore) removeModuleZip(mod coordinates.Module) error {
	sAtV := mod.AtVersion()

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	res, err := m.statements[deleteModuleZipSQL].ExecContext(ctx, sAtV)
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
		return errors.Errorf("expected exactly 1 row to be deleted for %s but got %d", mod, rowsAffected)
	}
	return nil
}

// for Index
func (m *mysqlStore) getModuleVersions(source string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := m.statements[selectModuleVersionsSQL].QueryContext(ctx, source)
	if err != nil {
		m.emitter.Count("db-select-module-versions-failure", 1)
		return nil, errors.Wrapf(err, "failed to query versions for %s", source)
	}

	defer ignoreClose(rows)

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

// for Index
func (m *mysqlStore) getVersionInfo(mod coordinates.Module) (repository.RevInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	row := m.statements[selectModuleVersionInfoSQL].QueryRowContext(ctx, mod.Source, mod.Version)
	var revInfo repository.RevInfo

	var contents []byte
	err := row.Scan(&contents)
	if err != nil {
		if err == sql.ErrNoRows {
			return revInfo, errors.New("module not in index")
		}
		m.emitter.Count("db-select-revinfo-failure", 1)
		return revInfo, errors.Wrapf(err, "failed to read row for sql: %+v", m.statements[selectModuleVersionInfoSQL])
	}

	err = json.Unmarshal(contents, &revInfo)
	return revInfo, err
}

// for Index
func (m *mysqlStore) getRegistryID(mod coordinates.Module) (bool, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	row := m.statements[selectRegistryIDSQL].QueryRowContext(ctx, mod.Source, mod.Version)

	var regid int64
	err := row.Scan(&regid)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0, nil
		}
		m.emitter.Count("db-select-regid-failure", 1)
		return false, 0, errors.Wrapf(err, "failed to read row for sql: %+v", m.statements[selectRegistryIDSQL])
	}

	return true, regid, nil
}

// for Index
func (m *mysqlStore) getGoMod(mod coordinates.Module) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	row := m.statements[selectGoModFileSQL].QueryRowContext(ctx, mod.Source, mod.Version)

	var contents []byte
	if err := row.Scan(&contents); err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("module not in index")
		}
		m.emitter.Count("db-select-gomod-failure", 1)
		return "", errors.Wrapf(err, "failed to read row for sql: %+v", m.statements[selectGoModFileSQL])
	}

	return string(contents), nil
}

// for Index
func (m *mysqlStore) ids() ([]int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := m.statements[selectAllRegistryIDsSQL].QueryContext(ctx)
	if err != nil {
		m.emitter.Count("db-select-ids-failure", 1)
		return nil, errors.Wrapf(err, "failed to query ids")
	}
	defer ignoreClose(rows)

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

// for Index
func (m *mysqlStore) countSourcesAndVersions() (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := m.statements[countVersionsSQL].QueryContext(ctx)
	if err != nil {
		m.emitter.Count("db-count-versions-failure", 1)
		return 0, 0, errors.Wrapf(err, "failed to query sources")
	}
	defer ignoreClose(rows)

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

func ignoreClose(c io.Closer) {
	_ = c.Close()
}
