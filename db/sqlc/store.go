package db

import (
	"context"
	"database/sql"
	"fmt"
)

// SQLStore provides all the functions to execute db queries and transections
type Store interface{
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore provides all the functions to execute SQL queries and transections
type SQLStore struct {
	*Queries
	//create a new db transections
	db *sql.DB
}

// NewStore creates a new Store

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transection
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {

	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	q := New(tx)

	if err = fn(q); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v,rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()

}

// TransferTxParams contains the input parameters of the transfer transection
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id`
	ToAccountID   int64 `json:"to_account_id`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transection
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

//TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, add account entries, add update account's balance within a single database transection

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		ctx.Value(txKey)
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Account:       arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.FromAccountID,
			Account:   -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.ToAccountID,
			Account:   +arg.Amount,
		})
		if err != nil {
			return err
		}

		//get account -> update it's balance
		// account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		// if err != nil {
		// 	return nil
		// }

		// account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		// if err != nil {
		// 	return nil
		// }

		// result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      arg.ToAccountID,
		// 	Balance: account2.Balance + arg.Amount,
		// })
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
			if err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}
		return nil
	})
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}
	return
}
