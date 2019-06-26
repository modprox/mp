package store

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"oss.indeed.com/go/modprox/pkg/coordinates"
	"oss.indeed.com/go/modprox/pkg/metrics/stats"
	"oss.indeed.com/go/modprox/pkg/repository"
)

type testSuite struct {
	suite.Suite

	ctx     context.Context
	db      *sql.DB
	subject *mysqlStore
}

func (s *testSuite) Test_ZipStore_PutZip_GetZip() {
	t := s.T()

	module := coordinates.Module{Source: "src1", Version: "v1.2.3"}
	blob := repository.Blob([]byte(string("hello")))
	err := s.subject.PutZip(module, blob)
	require.NoError(t, err)

	actual, err := s.subject.GetZip(module)
	require.NoError(t, err)
	require.Equal(t, blob, actual)
}

func (s *testSuite) Test_ZipStore_GetZip_NotFound() {
	t := s.T()

	module := coordinates.Module{Source: "src1", Version: "v1.2.3"}
	_, err := s.subject.GetZip(module)
	require.Error(t, err)
}

func (s *testSuite) Test_ZipStore_DelZip() {
	t := s.T()

	module := coordinates.Module{Source: "src1", Version: "v1.2.3"}
	blob := repository.Blob([]byte(string("hello")))
	err := s.subject.PutZip(module, blob)
	require.NoError(t, err)
	_, err = s.subject.GetZip(module)
	require.NoError(t, err)

	err = s.subject.DelZip(module)
	require.NoError(t, err)

	_, err = s.subject.GetZip(module)
	require.Error(t, err)
}

func (s *testSuite) Test_ZipStore_DelZip_NotFound() {
	t := s.T()

	module := coordinates.Module{Source: "src1", Version: "v1.2.3"}
	err := s.subject.DelZip(module)
	require.Error(t, err)
}

func (s *testSuite) Test_Index_PutGet() {
	t := s.T()

	module := coordinates.Module{Source: "src1", Version: "v1.2.3"}
	addition := ModuleAddition{Mod: module, UniqueID: int64(1234), ModFile: "foobar"}
	err := s.subject.Put(addition)
	require.NoError(t, err)

	actual, err := s.subject.Mod(module)
	require.NoError(t, err)
	require.Equal(t, actual, addition.ModFile)
}

func (s *testSuite) Test_Index_GetModNotFound() {
	t := s.T()

	module := coordinates.Module{Source: "src1", Version: "v1.2.3"}
	_, err := s.subject.Mod(module)
	require.Error(t, err)
	require.Equal(t, "module not in index", err.Error())
}

func (s *testSuite) Test_Index_Info() {
	t := s.T()

	module := coordinates.Module{Source: "src1", Version: "v1.2.3"}
	addition := ModuleAddition{Mod: module, UniqueID: int64(1234), ModFile: "foobar"}
	err := s.subject.Put(addition)
	require.NoError(t, err)

	actual, err := s.subject.Info(module)
	require.NoError(t, err)
	require.Equal(t, repository.RevInfo{Version: module.Version}, actual)
}

func (s *testSuite) Test_Index_InfoNotFound() {
	t := s.T()

	_, err := s.subject.Info(coordinates.Module{Source: "src1", Version: "v1.2.3"})
	require.Error(t, err)
	require.Equal(t, "module not in index", err.Error())
}

func (s *testSuite) Test_Index_Contains() {
	t := s.T()

	module := coordinates.Module{Source: "src1", Version: "v1.2.3"}
	addition := ModuleAddition{Mod: module, UniqueID: int64(1234), ModFile: "foobar"}
	err := s.subject.Put(addition)
	require.NoError(t, err)

	contains, id, err := s.subject.Contains(module)
	require.NoError(t, err)
	require.True(t, contains)
	require.Equal(t, addition.UniqueID, id)

	contains, id, err = s.subject.Contains(coordinates.Module{Source: "src1", Version: "v1.2.4"})
	require.NoError(t, err)
	require.False(t, contains)
}

func (s *testSuite) Test_Index_IDs() {
	t := s.T()

	ids := []int64{1, 2, 3, 4, 8, 11, 12, 13}
	for _, id := range ids {
		module := coordinates.Module{Source: "src1", Version: fmt.Sprintf("v1.2.%d", id)}
		addition := ModuleAddition{Mod: module, UniqueID: id, ModFile: "foobar"}
		err := s.subject.Put(addition)
		require.NoError(t, err)
	}

	actualIDs, err := s.subject.IDs()
	require.NoError(t, err)
	require.Equal(t, Ranges([][2]int64{{1, 4}, {8, 8}, {11, 13}}), actualIDs)
}

func (s *testSuite) Test_Index_UpdateID() {
	t := s.T()

	module := coordinates.Module{Source: "src1", Version: "v1.2.3"}
	addition := ModuleAddition{Mod: module, UniqueID: int64(1234), ModFile: "foobar"}
	err := s.subject.Put(addition)
	require.NoError(t, err)

	err = s.subject.UpdateID(coordinates.SerialModule{SerialID: 1235, Module: module})
	require.NoError(t, err)

	_, id, err := s.subject.Contains(module)
	require.NoError(t, err)
	require.Equal(t, int64(1235), id)
}

func (s *testSuite) Test_Index_Versions() {
	t := s.T()

	ids := []int64{1, 2, 3}
	source := "src1"
	for _, id := range ids {
		module := coordinates.Module{Source: source, Version: fmt.Sprintf("v1.2.%d", id)}
		addition := ModuleAddition{Mod: module, UniqueID: id, ModFile: "foobar"}
		err := s.subject.Put(addition)
		require.NoError(t, err)
	}

	actualVersions, err := s.subject.Versions(source)
	require.NoError(t, err)
	require.Equal(t, []string{"v1.2.1", "v1.2.2", "v1.2.3"}, actualVersions)

	actualVersions, err = s.subject.Versions("doesn't exist")
	require.NoError(t, err)
	require.Equal(t, []string{}, actualVersions)
}

func (s *testSuite) Test_Index_Remove() {
	t := s.T()

	module := coordinates.Module{Source: "src1", Version: "v1.2.3"}
	addition := ModuleAddition{Mod: module, UniqueID: int64(1234), ModFile: "foobar"}
	err := s.subject.Put(addition)
	require.NoError(t, err)

	exists, _, err := s.subject.Contains(module)
	require.NoError(t, err)
	require.True(t, exists)

	err = s.subject.Remove(module)
	require.NoError(t, err)

	exists, _, err = s.subject.Contains(module)
	require.NoError(t, err)
	require.False(t, exists)

	err = s.subject.Remove(module)
	require.NoError(t, err)
}

func (s *testSuite) Test_Index_Summary() {
	t := s.T()

	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			module := coordinates.Module{Source: fmt.Sprintf("src%d", i), Version: fmt.Sprintf("v1.2.%d", j)}
			addition := ModuleAddition{Mod: module, UniqueID: int64(i * j), ModFile: "foobar"}
			err := s.subject.Put(addition)
			require.NoError(t, err)
		}
	}

	totalSources, totalVersions, err := s.subject.Summary()
	require.NoError(t, err)
	require.Equal(t, 3, totalSources)
	require.Equal(t, 12, totalVersions)
}

func openTestDB(t *testing.T, forSchema bool) *sql.DB {
	config := mysql.Config{
		Net:                  "tcp",
		User:                 "docker",
		Passwd:               "docker",
		Addr:                 "localhost:3307",
		DBName:               "modproxdb-prox",
		AllowNativePasswords: true,
		ReadTimeout:          1 * time.Minute,
		WriteTimeout:         1 * time.Minute,
		MultiStatements:      forSchema,
	}
	if os.Getenv("TRAVIS") == "true" {
		config.Addr = "localhost:3306"
		config.User = "travis"
		config.Passwd = ""
		config.DBName = "modproxdb" // for some reason the mysql cli can't handle db names with dashes
	}
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}
	return db
}

func mustCloseDB(t *testing.T, db *sql.DB) {
	if err := db.Close(); err != nil {
		t.Fatalf("error closing DB: %v", err)
	}
}

func dropTables(t *testing.T, db *sql.DB) {
	tables := []string{
		"proxy_module_zips",
		"proxy_modules_index",
	}
	for _, table := range tables {
		if _, err := db.Exec("DROP TABLE IF EXISTS " + table); err != nil {
			t.Fatalf("error dropping table '%s': %v", table, err)
		}
	}
}

func createTables(t *testing.T, db *sql.DB) {
	schemaFile := filepath.Join("..", "..", "..", "..", "hack", "sql", "mysql-prox", "modproxdb.sql")
	schemaFileContents, err := ioutil.ReadFile(schemaFile)
	if err != nil {
		t.Fatalf("error reading schema file '%s': %v", schemaFile, err)
	}
	_, err = db.Exec(string(schemaFileContents))
	if err != nil {
		t.Fatalf("error applying schema file '%s': %v", schemaFile, err)
	}
}

func initSchema(t *testing.T) {
	db := openTestDB(t, true)
	defer mustCloseDB(t, db)

	dropTables(t, db)
	createTables(t, db)
}

type noopEmitter struct{}

var _ stats.Sender = (*noopEmitter)(nil)

func (*noopEmitter) Count(metric string, i int) {}

func (*noopEmitter) Gauge(metric string, n int) {}

func (*noopEmitter) GaugeMS(metric string, t time.Time) {}

func (s *testSuite) SetupTest() {
	t := s.T()

	initSchema(t)

	s.ctx = context.Background()
	s.db = openTestDB(t, false)
	tested, err := New("myql", s.db, &noopEmitter{})
	if err != nil {
		t.Fatalf("error connecting to db: %v", err)
	}
	s.subject = tested
}

func (s *testSuite) TearDownTest() {
	mustCloseDB(s.T(), s.db)
}

func TestMySQLSIDStore(t *testing.T) {
	suite.Run(t, new(testSuite))
}
