// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: transfar.sql

package sqlc

import (
	"context"
)

const createTransfar = `-- name: CreateTransfar :one
INSERT INTO transfars (
  amount,
  from_account_id,
  to_account_id
) VALUES (
  $1, $2, $3
) RETURNING id, amount, status, from_account_id, to_account_id, created_at
`

type CreateTransfarParams struct {
	Amount        int64 `json:"amount"`
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
}

func (q *Queries) CreateTransfar(ctx context.Context, arg CreateTransfarParams) (Transfar, error) {
	row := q.db.QueryRowContext(ctx, createTransfar, arg.Amount, arg.FromAccountID, arg.ToAccountID)
	var i Transfar
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.Status,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.CreatedAt,
	)
	return i, err
}

const deleteTransfar = `-- name: DeleteTransfar :exec
DELETE FROM transfars
WHERE id = $1
`

func (q *Queries) DeleteTransfar(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTransfar, id)
	return err
}

const getTransfar = `-- name: GetTransfar :one
SELECT id, amount, status, from_account_id, to_account_id, created_at FROM transfars
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTransfar(ctx context.Context, id int64) (Transfar, error) {
	row := q.db.QueryRowContext(ctx, getTransfar, id)
	var i Transfar
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.Status,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.CreatedAt,
	)
	return i, err
}

const listTransfars = `-- name: ListTransfars :many
SELECT id, amount, status, from_account_id, to_account_id, created_at FROM transfars
ORDER BY id DESC
LIMIT $1 OFFSET $2
`

type ListTransfarsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListTransfars(ctx context.Context, arg ListTransfarsParams) ([]Transfar, error) {
	rows, err := q.db.QueryContext(ctx, listTransfars, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transfar
	for rows.Next() {
		var i Transfar
		if err := rows.Scan(
			&i.ID,
			&i.Amount,
			&i.Status,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTransfar = `-- name: UpdateTransfar :one
UPDATE transfars
  SET status = $2
WHERE id = $1
RETURNING id, amount, status, from_account_id, to_account_id, created_at
`

type UpdateTransfarParams struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

func (q *Queries) UpdateTransfar(ctx context.Context, arg UpdateTransfarParams) (Transfar, error) {
	row := q.db.QueryRowContext(ctx, updateTransfar, arg.ID, arg.Status)
	var i Transfar
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.Status,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.CreatedAt,
	)
	return i, err
}
