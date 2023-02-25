package api

import (
	"context"
	"log"
	"net/http"
	"simplebank/db/sqlc"
	lib "simplebank/libs"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID    int64 `json:"from_account_id" binding:"required,min=1" `
	ToAccountID    int64 `json:"to_account_id" binding:"required,min=1" `
	Amount    int64 `json:"amount" binding:"required,gt=0" `
	Currency    string `json:"currency" binding:"required,currency" `
}


func (server *Server) transferAmount(ctx *gin.Context) {


	var req transferRequest


	if err := ctx.ShouldBindJSON(&req); err != nil {
		lib.HandleGinError(ctx, err)
		return 
	}


	
	arg := sqlc.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID: req.ToAccountID,
		Amount: req.Amount,
		
	}

	hasMatchingCurrency, err:= server.store.HasMatchingCurrency(ctx, sqlc.HasMatchingCurrencyParams{
		FromAccountID: arg.FromAccountID,
		ToAccountID: arg.ToAccountID,
		Currency: req.Currency,
	})

	log.Println("MATHICNG CURRENCY", hasMatchingCurrency, err)

	if lib.HasError(err) {
		lib.HandleGinErrorWithStaus(ctx,  http.StatusBadRequest ,err)
		return 
	}


    response, err := server.store.TransferTx(context.Background(), arg)

	if lib.HasError(err) {
		lib.HandleGinError(ctx, err)
		return
	}

	// response := sqlc.Account{}

	lib.HandleGinSuccess(ctx, response )


}


