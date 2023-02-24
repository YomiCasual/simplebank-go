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


func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)

	lib.HandleError(err);

	testQueries = sqlc.New(conn)

	os.Exit(m.Run()) 

}