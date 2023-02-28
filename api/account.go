package api

import (
	"context"
	"net/http"
	"simplebank/db/sqlc"
	lib "simplebank/libs"
	"strconv"

	"github.com/gin-gonic/gin"
)


type createAccountRequst struct {
	Owner    string `json:"owner" binding:"required" `
	Currency string `json:"currency" binding:"required,currency"`
}
type updateAccountRequest struct {
	Amount    int64 `json:"amount" binding:"required"`
}
type updateAccountRequestUri struct {
	ID    int64 `uri:"id" binding:"required,min=1"`
}
type listAccountsRequest struct {
	Page    int32 `form:"page" binding:"min=1"`
	PageSize    int32 `form:"pageSize" binding:"min=1"`
}


type getAccountRequest struct {
	ID    int64 `uri:"id" binding:"required,min=1"`
}

type deleteResponse interface {

}

func (server *Server) createAccount(ctx *gin.Context) {

	var req createAccountRequst;

	if err := ctx.ShouldBindJSON(&req); err !=nil {
		lib.HandleGinError(ctx, err)
		return 
	}

	arg := sqlc.CreateAccountParams{Balance: 0, Owner: req.Owner, Currency: req.Currency}


	account, err := server.store.CreateAccount(context.Background(), arg)


	lib.HandleAllErrors(ctx, err,  "Error")

	 ctx.JSON(http.StatusOK, account)

}
func (server *Server) listAccounts(ctx *gin.Context) {


	var params listAccountsRequest;

	if err := ctx.ShouldBindQuery(&params); err != nil {
		lib.HandleGinError(ctx, err)
		return 
	}

	arg := sqlc.ListAccountsParams{
		Limit: params.PageSize,
		Offset: (params.Page - 1) * params.PageSize,
	}


	accounts, err := server.store.ListAccounts(context.Background(), arg)

	if lib.HasError(err) {
		lib.HandleGinErrorWithStaus(ctx, http.StatusInternalServerError, err)
		return 
	}


	lib.HandleGinSuccess(ctx, accounts)

}
func (server *Server) updateAccount(ctx *gin.Context) {

	var reqUri updateAccountRequestUri

	var reqBody updateAccountRequest

	if err := ctx.ShouldBindUri(&reqUri); err !=nil {
		lib.HandleGinError(ctx, err)
		return 
	}


	if err := ctx.ShouldBindJSON(&reqBody); err !=nil {
		lib.HandleGinError(ctx, err)
		return 
	}

	arg := sqlc.UpdateAccountBalanceParams{
		Amount: reqBody.Amount,
		ID: reqUri.ID,
	}

	account, err := server.store.UpdateAccountBalance(context.Background(), arg)

	if lib.HasError(err) {
		lib.HandleGinErrorWithStaus(ctx, http.StatusInternalServerError, err)
		return 
	}

	lib.HandleGinSuccess(ctx, account)

}


func (server *Server) getAccount(ctx *gin.Context) {


	var req getAccountRequest


	if err := ctx.ShouldBindUri(&req); err !=nil {
		lib.HandleGinError(ctx, err)
		return 
	}

	account, err := server.store.GetAccount(context.Background(), req.ID)

	if lib.HasError(err) {
		if lib.IsSqlNotFounderror(err) {
			lib.HandleGinErrorWithStaus(ctx, http.StatusNotFound, err)
			return
		}
		lib.HandleGinErrorWithStaus(ctx, http.StatusInternalServerError, err)
		return 
	}

	lib.HandleGinSuccess(ctx, account)
}


func (server *Server) deleteAccount(ctx *gin.Context) {
	id := ctx.Param("id")


	result, err := strconv.ParseInt(id, 10, 64)

	if (lib.HasError(err)) {
		lib.HandleGinErrorWithStaus(ctx, http.StatusBadRequest, err)

	}

	newError := server.store.DeleteAccount(context.Background(), result)
	if lib.HasError(newError) {
		lib.HandleGinErrorWithStaus(ctx, http.StatusInternalServerError, newError)
		return 
	}

	var	account deleteResponse

	lib.HandleGinSuccess(ctx, account)
}