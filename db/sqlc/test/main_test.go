package test

import (
	"database/sql"
	"os"
	"simplebank/db/sqlc"
	lib "simplebank/libs"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://user:password@localhost:5432/simplebank?sslmode=disable"
)

var testQueries *sqlc.Queries

var testDb *sql.DB


func TestMain(m *testing.M) {

	var err error

	testDb, err = sql.Open(dbDriver, dbSource)

	lib.HandleError(err);

	testQueries = sqlc.New(testDb)

	os.Exit(m.Run()) 

}