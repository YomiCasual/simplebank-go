package main

import (
	"database/sql"
	"simplebank/api"
	"simplebank/db/sqlc"
	lib "simplebank/libs"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://user:password@localhost:5432/simplebank?sslmode=disable"
)

func main() {


	var err error

	conn, err := sql.Open(dbDriver, dbSource)

	lib.HandleError(err);
	
	server := api.NewServer(sqlc.NewStore(conn))

	err = server.Start("8080")

	lib.HandleError(err)


}