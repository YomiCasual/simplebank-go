package main

import (
	"database/sql"
	"simplebank/api"
	"simplebank/db/sqlc"
	lib "simplebank/libs"

	_ "github.com/lib/pq"
)



func main() {

	config, _ := lib.LoadConfig(".");


	

	var err error

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	lib.HandleError(err);
	
	server := api.NewServer(sqlc.NewStore(conn))

	err = server.Start(config.ServerAddress)

	lib.HandleError(err)


}