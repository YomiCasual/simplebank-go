package lib

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func HandleGinError(ctx *gin.Context, err error)  {
	 ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
}
func HandleGinErrorWithStaus(ctx *gin.Context,status int,  err error)  {
	 ctx.JSON(status, ErrorResponse(err))
}

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func HasError(err error) bool {
	return err != nil
}