package test

import (
	"context"
	"log"
	"simplebank/db/sqlc"
	"testing"

	"github.com/stretchr/testify/require"
)


func TestTransferTx(t *testing.T) {
	store := sqlc.NewStore(testDb)

	account1, _, _ := createAccount(t)
	account2, _, _ := createAccount(t)

	// run n concurrent transfer transactions

	n := 3;
	amount := int64(10)

	errs := make(chan error, n)
	results := make(chan sqlc.TransferTxResult, n)

	for i:=0; i < n;i++ {
		go func ()  {
			ctx := context.Background()
				result, err := store.TransferTx(ctx, sqlc.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID: account2.ID,
					Amount: int64(amount),
				})

				errs <- err
				results <- result
		}()
	}


	// close(errs)
	// close(results)

	// check results
	existed := make(map[int]bool)

	for i:=0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results


		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check balances
		log.Println("Account 1 balance", account1.Balance)
		log.Println("Account 2 balance", account2.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

	}



	// check the final updated balance
	updatedAccount1, err := store.GetAccountForUpdate(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccountForUpdate(context.Background(), account2.ID)
	require.NoError(t, err)


	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)

}