-- name: CreateTransfar :one
INSERT INTO transfars (
  amount,
  from_account_id,
  to_account_id
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetTransfar :one
SELECT * FROM transfars
WHERE id = $1 LIMIT 1;

-- name: ListTransfars :many
SELECT * FROM transfars
ORDER BY id DESC
LIMIT $1 OFFSET $2;

-- name: UpdateTransfar :one
UPDATE transfars
  SET status = $2
WHERE id = $1
RETURNING *;

-- name: DeleteTransfar :exec
DELETE FROM transfars
WHERE id = $1;