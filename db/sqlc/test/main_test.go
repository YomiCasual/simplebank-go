package test

import (
	"database/sql"
	"os"
	"simplebank/db/sqlc"
	lib "simplebank/libs"
	"testing"

	_ "github.com/lib/pq"
)



var testQueries *sqlc.Queries

var testDb *sql.DB


func TestMain(m *testing.M) {

	config, _ := lib.LoadConfig("./../../../");

	var err error

	testDb, err = sql.Open(config.DBDriver, config.DBSource)

	lib.HandleError(err);

	testQueries = sqlc.New(testDb)

	os.Exit(m.Run()) 

}