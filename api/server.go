package api

import (
	"fmt"
	"simplebank/db/sqlc"
	lib "simplebank/libs"
	"simplebank/token"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)


type TokenSelectorReturnType interface {
	map[string]token.Maker[int32]
}

type Server struct {
	config lib.Config
	store *sqlc.Store
	tokenMaker token.Maker[int32]
	router *gin.Engine
}

//   func TokenSelector[T TokenSelectorReturnType](secretKey string) T {

// 	paseToMaker, _ := token.NewPasetoMaker

// 	return  T{
// 		"paseto": paseToMaker,
// 		"jwt": token.NewJWTMaker(secretKey),
// 	}
// }



func NewServer( config lib.Config, store *sqlc.Store) (*Server, error) {


	tokenMaker, err := token.NewJWTMaker(config.SymmetricKey)

	if lib.HasError(err) {
		return nil, fmt.Errorf("cannot create token maker %d", err)
	}

	server := &Server{ config: config,  store: store, tokenMaker: tokenMaker }


	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}


	server.accountRoutes(&router.RouterGroup)
	server.userRoutes(&router.RouterGroup)
	server.transferRoutes(&router.RouterGroup)
	server.authRoutes(&router.RouterGroup)


	server.router = router

	return server, nil
}





func (server *Server) accountRoutes(router *gin.RouterGroup)  {
	accounts := router.Group("/accounts")
	{
		accounts.GET("/", server.listAccounts )
		accounts.POST("/", server.createAccount )
		accounts.PATCH("/:id", server.updateAccount)
		accounts.GET("/:id", server.getAccount)
		accounts.DELETE("/:id", server.deleteAccount)
	}
}

func (server *Server) userRoutes(router *gin.RouterGroup)  {
	user := router.Group("/users")
	{
		user.POST("/", server.createUser )
		user.GET("/", server.listUsers )
	}
}

func (server *Server) transferRoutes(router *gin.RouterGroup)  {
	//Users
	auth := router.Group("/auth")
	{
		auth.POST("/login", server.loginUser )
	}
}

func (server *Server) authRoutes(router *gin.RouterGroup)  {
	transfer := router.Group("/transfer")
	{
		transfer.POST("/", server.transferAmount )
	}
}


func (server *Server) Start(address string) error {
	return server.router.Run(address)
}



