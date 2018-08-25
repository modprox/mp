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
	ListSources() ([]repository.Module, error)
	Add([]repository.Module) (int, int, error)
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

type sourceTableRow struct {
	id      int
	source  string
	created int
}

type tagsTableRow struct {
	id       int
	tag      string
	created  int
	sourceID int
}

type scanRow struct {
	sourceTableRow
	tagsTableRow
}

func (s *store) ListSources() ([]repository.Module, error) {
	rows, err := s.statements[selectSourcesScanSQL].Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	modules := make([]repository.Module, 0, 10)
	for rows.Next() {
		var row scanRow
		if err := rows.Scan(
			&row.sourceTableRow.id,
			&row.sourceTableRow.source,
			&row.sourceTableRow.created,
			&row.tagsTableRow.id,
			&row.tagsTableRow.tag,
			&row.tagsTableRow.created,
			&row.tagsTableRow.sourceID,
		); err != nil {
			return nil, err
		}
		modules = append(modules, repository.Module{
			Source:  row.source,
			Version: row.tag,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return modules, nil
}

func (s *store) Add(modules []repository.Module) (int, int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()

	sourcesAdded := 0
	tagsAdded := 0
	for _, module := range modules {
		sourceID, addedSource, err := s.maybeAddSource(tx, module.Source)
		if err != nil {
			return 0, 0, err
		} else if addedSource {
			sourcesAdded++
		}

		addedTag, err := s.maybeAddTag(tx, sourceID, module.Version)
		if err != nil {
			return 0, 0, err
		} else if addedTag {
			tagsAdded++
		}
	}

	return sourcesAdded, tagsAdded, tx.Commit()
}

func (s *store) maybeAddTag(tx *sql.Tx, sourceID int64, version string) (bool, error) {
	result, err := tx.Stmt(s.statements[insertTagSQL]).Exec(version, sourceID)
	if err != nil {
		return false, err
	}
	return maybeAffectedN(result, 1)
}

func (s *store) maybeAddSource(tx *sql.Tx, source string) (int64, bool, error) {
	result, err := tx.Stmt(s.statements[insertSourceSQL]).Exec(source)
	if err != nil {
		return 0, false, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, false, err
	}
	added, err := maybeAffectedN(result, 1)
	return id, added, err
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
