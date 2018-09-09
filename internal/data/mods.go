package data

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/modprox/libmodprox/repository"
)

type moduleTR struct {
	id      int64
	source  string
	version string
	// created int
}

func (s *store) ListMods() ([]repository.ModInfo, error) {
	rows, err := s.statements[selectSourcesScanSQL].Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	modules := make([]repository.ModInfo, 0, 10)
	for rows.Next() {
		var row moduleTR
		if err := rows.Scan(
			&row.id,
			&row.source,
			&row.version,
			// &row.created,
		); err != nil {
			return nil, err
		}
		modules = append(modules, repository.ModInfo{
			Source:  row.source,
			Version: row.version,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return modules, nil
}

func (s *store) AddMods(modules []repository.ModInfo) (int, error) {
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

func (s *store) isModuleInDB(tx *sql.Tx, mod repository.ModInfo) (int64, bool, error) {
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

func (s *store) insertModuleInDB(tx *sql.Tx, mod repository.ModInfo) error {
	// the PQ library DOES NOT SUPPORT LastInsertId, DO NOT USE IT
	_, err := tx.Stmt(s.statements[insertModuleSQL]).Exec(
		mod.Source,
		mod.Version,
	)
	return err
}
