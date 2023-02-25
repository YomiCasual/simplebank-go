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
	accounts := router.Group("/accounts")
	{
		accounts.GET("/", server.listAccounts )
		accounts.POST("/", server.createAccount )
		accounts.PATCH("/:id", server.updateAccount)
		accounts.GET("/:id", server.getAccount)
		accounts.DELETE("/:id", server.deleteAccount)
	}


	server.router = router

	return server
}


func (server *Server) Start(address string) error {
	return server.router.Run(address)
}



