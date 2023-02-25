package lib

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)


func HandleGinError(ctx *gin.Context, err error)  {
	 ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
}
func HandleGinErrorWithStaus(ctx *gin.Context,status int,  err error)  {
	 ctx.JSON(status, ErrorResponse(err))
}
func HandleGinErrorWithStatusAndMessage(ctx *gin.Context,status int,  message string)  {
	 ctx.JSON(status, gin.H{"success": false, "message": message})
}

func ErrorResponse(err error) gin.H {
	return gin.H{ "success": false, "error": err.Error()}
}

func HasError(err error) bool {
	return err != nil
}
func IsSqlNotFounderror(err error) bool {
	return err != sql.ErrNoRows
}

func HandleGinSuccess(ctx *gin.Context, response interface{})  {
	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
		"success": true,
	})
}