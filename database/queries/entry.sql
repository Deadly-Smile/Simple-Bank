-- name: CreateEntry :one
INSERT INTO entries (
  amount,
  account_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
ORDER BY id DESC
LIMIT $1 OFFSET $2;

-- -- name: UpdateEntry :one
-- UPDATE entries
--   SET username = $2
-- WHERE id = $1
-- RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;