package data

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"

	"github.com/pkg/errors"

	"oss.indeed.com/go/modprox/pkg/coordinates"
)

type moduleTR struct {
	id      int64
	source  string
	version string
}

func (s *store) ListModules() ([]coordinates.SerialModule, error) {
	start := time.Now()
	rows, err := s.statements[selectSourcesScanSQL].Query()
	if err != nil {
		s.emitter.Count("db-list-modules-failure", 1)
		return nil, err
	}
	defer rows.Close()

	mods, err := modulesFromRows(rows)
	if err != nil {
		s.emitter.Count("db-list-modules-failure", 1)
		return nil, err
	}

	s.emitter.GaugeMS("db-list-modules-elapsed-ms", start)
	return mods, nil
}

func listOfIDs(ids []int64) string {
	s := make([]string, 0, len(ids))
	for _, id := range ids {
		s = append(s, fmt.Sprintf("%d", id))
	}
	return strings.Join(s, ",")
}

func (s *store) ListModulesByIDs(ids []int64) ([]coordinates.SerialModule, error) {
	start := time.Now()
	mods, err := s.listModulesByIDs(ids)
	if err != nil {
		s.emitter.Count("db-list-modules-by-id-failure", 1)
		return nil, err
	}

	s.emitter.GaugeMS("db-list-modules-by-id-elapsed-ms", start)
	return mods, nil
}

func (s *store) listModulesByIDs(ids []int64) ([]coordinates.SerialModule, error) {
	var rows *sql.Rows
	var err error

	if stmt, exists := s.statements[selectModulesByIDsSQL]; exists {
		rows, err = stmt.Query(
			pq.Array(ids),
		)
	} else {
		// do not execute the generated query with no elements in the list
		// because it just breaks
		if len(ids) == 0 {
			return nil, nil
		}

		// generate this query by hand for mysql, who's driver still doesn't know
		// what an argument of list is in 2018
		text := "select id, source, version from modules where id in (%s) order by id asc"
		q := fmt.Sprintf(text, listOfIDs(ids))
		rows, err = s.db.Query(q)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return modulesFromRows(rows)
}

func (s *store) ListModulesBySource(source string) ([]coordinates.SerialModule, error) {
	start := time.Now()
	mods, err := s.listModulesBySource(source)
	if err != nil {
		s.emitter.Count("db-list-modules-by-source-failure", 1)
		return nil, err
	}

	s.emitter.GaugeMS("db-list-modules-by-source-elapsed-ms", start)
	return mods, nil
}

func (s *store) listModulesBySource(source string) ([]coordinates.SerialModule, error) {
	rows, err := s.statements[selectModulesBySource].Query(source)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mods, err := modulesFromRows(rows)
	if err != nil {
		return nil, err
	}

	return mods, nil
}

func modulesFromRows(rows *sql.Rows) ([]coordinates.SerialModule, error) {
	modules := make([]coordinates.SerialModule, 0, 10)
	for rows.Next() {
		var row moduleTR
		if err := rows.Scan(
			&row.id,
			&row.source,
			&row.version,
		); err != nil {
			return nil, err
		}
		modules = append(modules, coordinates.SerialModule{
			SerialID: row.id,
			Module: coordinates.Module{
				Source:  row.source,
				Version: row.version,
			},
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return modules, nil
}

func (s *store) ListModuleIDs() ([]int64, error) {
	start := time.Now()
	ids, err := s.listModuleIDs()
	if err != nil {
		s.emitter.Count("db-list-module-ids-failure", 1)
		return nil, err
	}

	s.emitter.GaugeMS("db-list-module-ids-elapsed-ms", start)
	return ids, nil
}

func (s *store) listModuleIDs() ([]int64, error) {
	rows, err := s.statements[selectModuleIDScanSQL].Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]int64, 0, 1024)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (s *store) InsertModules(modules []coordinates.Module) (int, error) {
	start := time.Now()
	i, err := s.insertModules(modules)
	if err != nil {
		s.emitter.Count("db-insert-modules-failure", 1)
		return 0, err
	}

	s.emitter.GaugeMS("db-insert-modules-elapsed-ms", start)
	return i, nil
}

func (s *store) insertModules(modules []coordinates.Module) (int, error) {
	modulesAdded := 0

	for _, mod := range modules {
		s.log.Tracef("inserting module: %s", mod)

		// start a transaction for each module insert, because this logic
		// is more complicated than one might think in the efforts to preserve
		// very sequential ids (but still not guaranteed)
		tx, err := s.db.Begin()
		if err != nil {
			return 0, err
		}

		// does the module already exist in the db?
		_, exists, err := s.isModuleInDB(tx, mod)
		if err != nil {
			_ = tx.Rollback()
			s.log.Errorf("failed to check if module in db")
			return 0, err
		}

		// if not, add the module into the db
		if !exists {
			if err := s.insertModuleInDB(tx, mod); err != nil {
				_ = tx.Rollback()
				s.log.Errorf("failed to insert module into db")
				return 0, err
			}
			modulesAdded++
		}

		// end the transaction for this module
		if err := tx.Commit(); err != nil {
			_ = tx.Rollback()
			s.log.Errorf("failed to commit insert of module into db")
			return 0, err
		}
	}

	return modulesAdded, nil
}

func (s *store) isModuleInDB(tx *sql.Tx, mod coordinates.Module) (int64, bool, error) {
	rows, err := tx.Stmt(s.statements[selectModuleIDSQL]).Query(
		mod.Source,
		mod.Version,
	)
	if err != nil {
		return 0, false, err
	}
	defer rows.Close()
	return maybeGetID(rows)
}

func maybeGetID(rows *sql.Rows) (int64, bool, error) {
	var ids []int64 // expect 0 or 1
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return 0, false, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return 0, false, err
	}

	numIDs := len(ids)
	switch numIDs {
	case 0:
		return 0, false, nil
	case 1:
		return ids[0], true, nil
	default:
		return 0, false, errors.Errorf("expected 0 or 1 rows, got %d", numIDs)
	}
}

func (s *store) insertModuleInDB(tx *sql.Tx, mod coordinates.Module) error {
	// the PQ library DOES NOT SUPPORT LastInsertId, DO NOT USE IT
	_, err := tx.Stmt(s.statements[insertModuleSQL]).Exec(
		mod.Source,
		mod.Version,
	)
	return err
}

func (s *store) DeleteModuleByID(id int) error {
	_, err := s.statements[deleteModuleByIDSQL].Exec(id)
	return err
}
