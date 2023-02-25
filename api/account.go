package api

import (
	"context"
	"net/http"
	"simplebank/db/sqlc"
	lib "simplebank/libs"

	"github.com/gin-gonic/gin"
)


type createAccountRequst struct {
	Owner    string `json:"owner" binding:"required" `
	Currency string `json:"currency" binding:"required,oneof= EUR USD CAD"`
}

func (server *Server) createAccount(ctx *gin.Context) {

	var req createAccountRequst;

	if err := ctx.ShouldBindJSON(&req); err !=nil {
		lib.HandleGinError(ctx, err)
		return 
	}

	arg := sqlc.CreateAccountParams{Balance: 0, Owner: req.Owner, Currency: req.Currency}


	account, err := server.store.CreateAccount(context.Background(), arg)

	if lib.HasError(err) {
		lib.HandleGinErrorWithStaus(ctx, http.StatusInternalServerError, err)
		return 
	}


	 ctx.JSON(http.StatusOK, account)

}