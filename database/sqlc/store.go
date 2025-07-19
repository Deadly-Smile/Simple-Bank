package sqlc

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			fmt.Printf("store execTx error: %v, rollback error: %v\n", err, rbErr)
			return rbErr
		}
		return err
	}

	return tx.Commit()
}

type TransfarParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransfarResult struct {
	Transfar    Transfar `json:"transfar"`
	ToEntry     Entry    `json:"to_entry"`
	FromEntry   Entry    `json:"from_entry"`
	ToAccount   Account  `json:"to_account"`
	FromAccount Account  `json:"from_account"`
}

func (store *Store) TransfarTx(ctx context.Context, transfarPrems TransfarParams) (TransfarResult, error) {
	var result TransfarResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfar, err = q.CreateTransfar(ctx, CreateTransfarParams{
			Amount:        transfarPrems.Amount,
			FromAccountID: transfarPrems.FromAccountID,
			ToAccountID:   transfarPrems.ToAccountID,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: transfarPrems.FromAccountID,
			Amount:    -transfarPrems.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: transfarPrems.ToAccountID,
			Amount:    transfarPrems.Amount,
		})
		if err != nil {
			return err
		}

		if transfarPrems.FromAccountID < transfarPrems.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, transfarPrems.FromAccountID, transfarPrems.ToAccountID, -transfarPrems.Amount, transfarPrems.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, transfarPrems.ToAccountID, transfarPrems.FromAccountID, transfarPrems.Amount, -transfarPrems.Amount)
		}

		return err
	})
	return result, err
}

func addMoney(ctx context.Context, q *Queries, accountID1 int64, accountID2 int64, amount1 int64, amount2 int64) (account1 Account, account2 Account, err error) {
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
