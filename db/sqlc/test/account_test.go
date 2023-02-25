package test

import (
	"context"
	"database/sql"
	"simplebank/db/sqlc"
	lib "simplebank/libs"
	"testing"

	"github.com/stretchr/testify/require"
)


func TestCreatingAccount(t *testing.T) {
	createAccount(t)
}

func TestGetAccount(t *testing.T) {

	account,_, _ := createAccount(t)

	
	userAccount, err := testQueries.GetAccount(context.Background(), account.ID)


	require.NoError(t, err)
	require.NotEmpty(t, userAccount)
	require.NotZero(t, userAccount.ID)
	require.NotZero(t, userAccount.CreatedAt)
	
	require.Equal(t, account.Owner, userAccount.Owner)
	require.Equal(t, account.Currency, userAccount.Currency)

}


func TestUpdateAccount(t *testing.T) {

	account, _, _ := createAccount(t)

	arg := sqlc.UpdateAccountParams{
		ID: account.ID,
		Balance: lib.RandomMoney(),
	}

	updatedAccount, err := testQueries.UpdateAccount(context.Background(), arg)


	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)
	require.NotZero(t, updatedAccount.ID)
	require.NotZero(t, updatedAccount.CreatedAt)

	require.Equal(t, account.Owner, updatedAccount.Owner)
	require.NotEqual(t, account.Balance, updatedAccount.Balance)

}


func TestDeleteAccount(t *testing.T) {

	account1,_, _ := createAccount(t)

	
	 err := testQueries.DeleteAccount(context.Background(), account1.ID)

	 require.NoError(t, err)

	 account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)

}
func TestListAccount(t *testing.T) {

	 createAccount(t)

	
	 accounts, err := testQueries.ListAccounts(context.Background())

	 require.NoError(t, err)


	require.NoError(t, err)
	require.NotEmpty(t, accounts)

}


func createAccount(t *testing.T) (sqlc.Account, sqlc.CreateAccountParams, error ) {
	arg := sqlc.CreateAccountParams{
		Owner: lib.RandomOwner(),
		Balance: lib.RandomMoney(),
		Currency: lib.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)



	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account,arg,  err

}