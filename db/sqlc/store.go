package sqlc

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	queries *Queries
	db      *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		queries: New(db),
	}
}

// execTx. if any transaction happen in this function where to fail, it all will be rollback to its initial state
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	n := New(tx)

	fnErr := fn(n)
	if fnErr != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return fmt.Errorf("this is transaction error %v an func error %v", txErr, err)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs money transfer from one account to another
// it  creates a transfer record, add account entries, and update accounts balance within a single database transaction
func (store *Store) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	var err error

	var TransferTxFn = func(q *Queries) error {
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccountId,
			ToAccountID:   args.ToAccountId,
			Amount:        args.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountId,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountId,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}

		acc1, err := q.GetAccount(ctx, args.FromAccountId)
		if err != nil {
			return err
		}

		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      args.FromAccountId,
			Balance: acc1.Balance - args.Amount,
		})
		if err != nil {
			return err
		}

		acc2, err := q.GetAccount(ctx, args.ToAccountId)
		if err != nil {
			return err
		}

		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      args.ToAccountId,
			Balance: acc2.Balance - args.Amount,
		})
		if err != nil {
			return err
		}

		return nil

	}
	store.execTx(ctx, TransferTxFn)

	return result, nil

}
