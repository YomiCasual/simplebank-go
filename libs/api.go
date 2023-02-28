package lib

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)


func HandleGinError(ctx *gin.Context, err error)  {
	 ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
}
func HandleGinErrorWithStaus(ctx *gin.Context,status int,  err error)  {
	 ctx.JSON(status, ErrorResponse(err))
}
func HandleGinErrorWithStatusAndMessage(ctx *gin.Context,status int,  message string)  {
	 ctx.JSON(status, gin.H{"success":false , "message": message,})
}
func HandleGinErrorWithStatusAndMessageWithError(ctx *gin.Context,status int, err error,  message string)  {
	 ctx.JSON(status, gin.H{"success":false , "error": err.Error(), "message": message,})
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

func HandleAllErrors(ctx *gin.Context, err error,  message string) {

		pqErr, ok := err.(*pq.Error)
		
		if (ok) {
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				HandleGinErrorWithStatusAndMessageWithError(ctx, http.StatusInternalServerError,err,  message);
				return
				case "unique_violation": 
				HandleGinErrorWithStatusAndMessageWithError(ctx, http.StatusInternalServerError, err, message);
				return

			}
		}

		HandleGinErrorWithStaus(ctx, http.StatusInternalServerError, err)

}