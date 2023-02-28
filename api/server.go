package api

import (
	"simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)


type Server struct {
	store *sqlc.Store
	router *gin.Engine
}


func NewServer(store *sqlc.Store) *Server {
	server := &Server{ store: store }


	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

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

	//Accounts
	transfer := router.Group("/transfer")
	{
		transfer.POST("/", server.transferAmount )
	}


	//Users
	user := router.Group("/users")
	{
		user.POST("/", server.createUser )
		user.GET("/", server.listUsers )
	}


	server.router = router

	return server
}


func (server *Server) Start(address string) error {
	return server.router.Run(address)
}



