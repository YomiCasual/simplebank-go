package api

import (
	"context"
	"net/http"
	"simplebank/db/sqlc"
	lib "simplebank/libs"
	"time"

	"github.com/gin-gonic/gin"
)




type createUserRequest struct {
	Username    string `json:"username" binding:"required,alphanum" `
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}
type UserResponse struct {
	ID                int64     `json:"id"`
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	CreatedAt         time.Time `json:"createdAt"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
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


	hashed_password, err := lib.PasswordCrypt().HashPassword(req.Password);

	if lib.HasError(err) {
		lib.HandleGinErrorWithStaus(ctx, http.StatusInternalServerError, err)
		return
	}

	arg := sqlc.CreateUserParams{
	Username: req.Username,
	HashedPassword: hashed_password,
	FullName: req.FullName,
	Email: req.Email,
	}


	user, err := server.store.CreateUser(context.Background(), arg)


	
	if lib.HasError(err) {
		lib.HandleAllErrors(ctx, err,  "User with this details exist")
		return
	}
	
	newUser := UserResponse{
		Username: user.Username,
		ID: user.ID,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		FullName: user.FullName,
		PasswordChangedAt: user.PasswordChangedAt,
	}

	lib.HandleGinSuccess(ctx, newUser)
}

func (server *Server) listUsers(ctx *gin.Context) {


	var params listUserRequest;

	
	if err := ctx.BindQuery(&params); err != nil {

		if params.Page == 0 || params.PageSize == 0 {
			params.Page = 1
			params.PageSize = 5
		}
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


