package api

import (
	"errors"
	"net/http"
	lib "simplebank/libs"
	"simplebank/token"
	"strings"

	"github.com/gin-gonic/gin"
)


const (
	authorizationHeaderKey = "authorization"
	authorizationBearerType = "bearer"
	authorizationUser = "authorization_user"
)


func (server *Server) authMiddleware() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if (len(authorizationHeader) == 0) {
			returnMiddleWareError(ctx, "authorization header is not provided")
			return 
		}

		fields := strings.Fields(authorizationHeader)

		if (len(fields) < 2 || strings.ToLower(fields[0]) != authorizationBearerType ) {
			returnMiddleWareError(ctx, "invalid authorization type")
			return
		}

		accessToken := fields[1]

		payload, err := server.tokenMaker.VerifyToken(accessToken)

		if lib.HasError(err) {
			returnMiddleWareError(ctx, "invalid token")
			return
		}

		ctx.Set(authorizationUser, payload)
		ctx.Next()
	}
}

func (server *Server) AuthUser(ctx *gin.Context) (*token.Payload, error ){
	authUser, ok := ctx.Get(authorizationUser)
	
	if !ok {
		return  nil, errors.New("invalid auth user")
	}
	
	assertedUser, ok := authUser.(*token.Payload)
	
	if !ok {
		return  nil, errors.New("Invalid auth user")
	}

	return assertedUser, nil
}

func returnMiddleWareError(ctx *gin.Context, msg string) {
	err := errors.New(msg)
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, lib.ErrorResponse(err))
}