package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		createdTransferId, _ := q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		transferId, _ := createdTransferId.LastInsertId()
		result.Transfer, err = q.GetTransfer(ctx, transferId)

		//result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParam{
		//	FromAcocuntID: arg.FromAccountID,
		//	ToAccountID:   arg.ToAccountID,
		//	Amount:        arg.Amount,
		//})
		if err != nil {
			return err
		}

		createdFromEntryId, _ := q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		fromEntryId, _ := createdFromEntryId.LastInsertId()
		result.FromEntry, err = q.GetEntry(ctx, fromEntryId)
		/*result.FromEntry, err := q.createEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})*/
		if err != nil {
			return err
		}

		createdToEntryId, _ := q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		toEntryId, _ := createdToEntryId.LastInsertId()
		result.ToEntry, err = q.GetEntry(ctx, toEntryId)

		/*result.ToEntry, err := q.createEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})*/
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
