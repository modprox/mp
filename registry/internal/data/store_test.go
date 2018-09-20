package data

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/modprox/mp/registry/config"
	"github.com/stretchr/testify/require"
)

var (
	postgresDSN = config.DSN{
		Address:  "127.0.0.1:5432",
		User:     "docker",
		Password: "docker",
		Database: "postgres",
	}
)

func createDB(t *testing.T, dbType string) (config.DSN, string) {
	if dbType == "mysql" {
		t.Fail() //todo: implement
	}
	return createPostgreSQL(t)
}

func cleanupDB(t *testing.T, kind string, dsn config.DSN) {
	if kind == "postgres" {
		db, err := connectPostgreSQL(dsn)
		require.NoError(t, err)

		_, err = db.Exec(fmt.Sprintf("drop database %s", dsn.Database))
		require.NoError(t, err)
	} else {
		// mysql
	}
}

//
//func createMySQL(t *testing.T) *sql.DB {
//	bs, err := ioutil.ReadFile("../../../hack/sql/mysql/modproxdb.sql")
//	require.NoError(t, err)
//
//	t.Log("sql:", string(bs))
//	return nil
//}

func createPostgreSQL(t *testing.T) (config.DSN, string) {
	dsn := config.DSN{
		Address:  "127.0.0.1:5432",
		User:     "docker",
		Password: "docker",
		Database: "postgres",
	} // used temporarily to create another test database

	db, err := connectPostgreSQL(dsn)
	require.NoError(t, err)

	// create temp test database and use that
	dbName := randomName()
	_, err = db.Exec(fmt.Sprintf("create database %s", dbName))
	require.NoError(t, err)

	// close postgres database connection
	err = db.Close()
	require.NoError(t, err)

	testDSN := config.DSN{
		Address:  "127.0.0.1:5432",
		User:     "docker",
		Password: "docker",
		Database: dbName,
	}

	db, err = connectPostgreSQL(testDSN)
	require.NoError(t, err)

	createTables(t, db, "../../../hack/sql/postgres/modproxdb.sql")
	err = db.Close()
	require.NoError(t, err)

	return testDSN, dbName
}

func randomName() string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	a := r.Int() % 10
	b := r.Int() % 10
	c := r.Int() % 10
	d := r.Int() % 10
	e := r.Int() % 10
	return fmt.Sprintf("mp_%d%d%d%d%d", a, b, c, d, e)
}

func createTables(t *testing.T, db *sql.DB, file string) {
	bs, err := ioutil.ReadFile(file)
	require.NoError(t, err)

	statements := strings.Split(string(bs), ";")
	for _, statement := range statements {
		_, err = db.Exec(statement)
		require.NoError(t, err)
	}
}

var dbTypes = []string{
	//	"mysql",
	"postgres",
}

func Test_Create(t *testing.T) {
	for _, dbType := range dbTypes {
		testCreate(t, dbType)
	}
}

func testCreate(t *testing.T, dbType string) {
	dsn, dbName := createDB(t, dbType)
	defer cleanupDB(t, dbType, dsn)

	t.Log("creating a database of dbType:", dbType, ", name:", dbName, "dsn:", dsn)
	store, err := Connect(dbType, dsn)
	require.NoError(t, err)
	_ = store
}
