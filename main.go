package main

import (
	"database/sql"
	"fmt"
	"log"
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
	
	server, err := api.NewServer(config, sqlc.NewStore(conn))

	if lib.HasError(err) {
		fmt.Println("error from server", err)
		log.Fatal("cannot get config")
	}

	err = server.Start(config.ServerAddress)

	lib.HandleError(err)


}