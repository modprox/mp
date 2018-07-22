package repositories

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"github.com/modprox/modprox-registry/internal/repositories/repository"
)

func Connect(config mysql.Config) (*sql.DB, error) {
	dsn := config.FormatDSN()
	return sql.Open("mysql", dsn)
}

type Store interface {
	List() ([]repository.Module, error)
	Append([]repository.Module) (int, error)
}

func New(db *sql.DB) (Store, error) {
	statements, err := load(db)
	if err != nil {
		return nil, err
	}
	return &store{
		db:         db,
		statements: statements,
	}, nil
}

type store struct {
	db         *sql.DB
	statements statements
}

func (s *store) List() ([]repository.Module, error) {
	rows, err := s.statements[selectAllRegistryInfos].Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	infos := make([]repository.Module, 0, 10)
	for rows.Next() {
		var origin string
		if err := rows.Scan(&origin); err != nil {
			return nil, err
		}
		infos = append(infos, repository.Module{Source: origin})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return infos, nil
}

func (s *store) Append(infos []repository.Module) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	numAdded := 0
	for _, info := range infos {
		if added, err := s.maybeAdd(tx, info); err != nil {
			return 0, err
		} else if added {
			numAdded++
		}
	}

	return numAdded, tx.Commit()
}

func (s *store) maybeAdd(tx *sql.Tx, info repository.Module) (bool, error) {
	result, err := tx.Stmt(s.statements[insertRegistryInfo]).Exec(info.Source, info.Version)
	if err != nil {
		return false, err
	}
	return maybeAffectedN(result, 1)
}

func maybeAffectedN(result sql.Result, n int) (bool, error) {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsAffected == 0 {
		return false, nil
	}

	if rowsAffected == int64(n) {
		return true, nil
	}

	return false, errors.Errorf("expected to affect %d rows, actually affected %d", n, rowsAffected)
}
