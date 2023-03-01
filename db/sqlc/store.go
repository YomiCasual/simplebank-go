package sqlc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
		Queries: New(db),
	}
}


var TxKey = struct{}{}


func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {

	tx, err := store.db.BeginTx(ctx, nil)
	
	if err != nil {
		return err
	}
	
	q := New(tx)
	
	err = fn(q)
	
	if (err != nil) {
	
		if rbEr := tx.Rollback(); rbEr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbEr)
		}

		return err
	}

	return tx.Commit()
}





type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}
type HasMatchingCurrencyParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Currency    string `json:"currency" binding:"required,currency" `
}
type IsOwnerAccountParams struct {
	FromAccountID int64 `json:"from_account_id"`
	AuthUserUsername string `json:"auth_user_username"`
}


type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_acount"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {

		var err error


		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})

		if (err != nil) {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})

		if (err != nil) {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})


		if (err != nil) {
			return err
		}

		//TODO: UPDATE ACCOUNTS BALANCE


		result.FromAccount, err = store.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
			ID: arg.FromAccountID,
			Amount:  - arg.Amount,
		})

		if (err != nil) {
			return errors.New("insufficent balance")
		}


		
		result.ToAccount, err = store.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
			ID: arg.ToAccountID, 
			Amount:  arg.Amount,
		})

		if (err != nil) {
			return err
		}


		return nil
	} )

	return result, err
}

func (store *Store) HasMatchingCurrency(ctx context.Context, arg HasMatchingCurrencyParams) (bool, error) {


 	hasMatchingAccount := false
	 
	 err := store.execTx(ctx, func(q *Queries) error {


		account1, err := q.GetAccount(ctx, arg.FromAccountID)

		if (err != nil) {
			err = errors.New("cannot get account with from_account_id")
			  return err
		}

		account2, err := q.GetAccount(ctx, arg.ToAccountID)

		if (err != nil) {
			err = errors.New("cannot get account with to_account_id")
			return err
		}

		if account1.Currency != account2.Currency || account1.Currency != arg.Currency {
			err := fmt.Errorf("currency mismatch of account1 [%v] and account 2 [%v] and transfer currency [%v]", account1.Currency, account2.Currency, arg.Currency)
			return err
		}

		return err
	} )


	return hasMatchingAccount, err
}


// func (store *Store) IsOwnerAccount(ctx context.Context, arg IsOwnerAccountParams) (bool, error) {


//  	isOwnerAccount := false
	 
// 	account, err := store.GetAccount(ctx, arg.FromAccountID);

// 	if lib.HasError(err) {
// 		accountError := fmt.Errorf("error getting account with id : %v", arg.FromAccountID)
// 		return isOwnerAccount, accountError
// 	}

// 	if account.Owner != arg.AuthUserUsername {
// 		accountError := fmt.Errorf("You don't have permission  to send from this account: %v", arg.FromAccountID)
// 		return isOwnerAccount, accountError
// 	}
	

// 	return hasMatchingAccount, err
// }

