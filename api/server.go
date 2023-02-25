package api

import (
	"simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
)


type Server struct {
	store *sqlc.Store
	router *gin.Engine
}


func NewServer(store *sqlc.Store) *Server {
	server := &Server{ store: store }

	router := gin.Default()

	//routes
	//Accounts

	router.POST("/accounts", server.createAccount )

	server.router = router

	return server
}


func (server *Server) Start(address string) error {
	return server.router.Run("localhost:" + address)
}



