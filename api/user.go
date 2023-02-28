package api

import (
	"context"
	"net/http"
	"simplebank/db/sqlc"
	lib "simplebank/libs"

	"github.com/gin-gonic/gin"
)




type createUserRequest struct {
	Username    string `json:"owner" binding:"required" `
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type listUserRequest struct {
	Page    int32 `form:"page" binding:"min=1"`
	PageSize    int32 `form:"pageSize" binding:"min=1"`
}




func (server *Server) createUser(ctx *gin.Context) {

	var req createUserRequest;

	if err := ctx.ShouldBindJSON(&req); err !=nil {
		lib.HandleGinError(ctx, err)
		return 
	}

	arg := sqlc.CreateUserParams{
	Username: req.Username,
	HashedPassword: req.Password,
	FullName: req.FullName,
	Email: req.Email,
	}


	user, err := server.store.CreateUser(context.Background(), arg)

	if lib.HasError(err) {
		lib.HandleGinErrorWithStaus(ctx, http.StatusInternalServerError, err)
		return 
	}

	lib.HandleGinSuccess(ctx, user)
}

func (server *Server) listUsers(ctx *gin.Context) {


	var params listUserRequest;

	if err := ctx.ShouldBindQuery(&params); err != nil {
		lib.HandleGinError(ctx, err)
		return 
	}

	arg := sqlc.ListUsersParams{
		Limit: params.PageSize,
		Offset: (params.Page - 1) * params.PageSize,
	}

	users, err := server.store.ListUsers(context.Background(), arg)

	if lib.HasError(err) {
		lib.HandleGinErrorWithStaus(ctx, http.StatusInternalServerError, err)
		return 
	}

	lib.HandleGinSuccess(ctx, users)

}


