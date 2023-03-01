package api

import (
	"context"
	"errors"
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
type loginRequest struct {
	UsernameEmail    string `json:"username_email" binding:"required" `
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User UserResponse `json:"user"`
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


func newUserResponse (user sqlc.User) UserResponse {
	return UserResponse{
		Username: user.Username,
		ID: user.ID,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		FullName: user.FullName,
		PasswordChangedAt: user.PasswordChangedAt,
	}
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
	
	newUser := newUserResponse(user)

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



func (server *Server) loginUser(ctx *gin.Context) {


	var params loginRequest;
	var user sqlc.User

	
	if err := ctx.ShouldBindJSON(&params); err != nil {
		lib.HandleGinErrorWithStaus(ctx, http.StatusInternalServerError, err)
		return
	}


	user, err := server.store.GetUserByUsername(context.Background(), params.UsernameEmail)

	if lib.HasError(err) {
		
		user, err = server.store.GetUserByEmail(context.Background(), params.UsernameEmail)

		if lib.HasError(err) {
			lib.HandleGinErrorWithStaus(ctx, http.StatusInternalServerError, errors.New("Invalid credentials"))
			return 
		}
	}
	
	hasValidPassword := lib.PasswordCrypt().CheckPassword(user.HashedPassword, params.Password)

	if (!hasValidPassword) {
		lib.HandleGinErrorWithStaus(ctx, http.StatusInternalServerError, errors.New("Invalid credentials"))
		return 
	}

	token, err := server.tokenMaker.CreateToken(user.Username, int32(user.ID), server.config.AccessTokenDuration );


	if lib.HasError(err) {
		lib.HandleGinErrorWithStaus(ctx, http.StatusInternalServerError, err)
		return 
	}

		
	newUser := newUserResponse(user)


	loginResponse := LoginResponse{
		Token: token,
		User: newUser,
	}


	lib.HandleGinSuccess(ctx, loginResponse)

}


