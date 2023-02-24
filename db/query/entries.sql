-- name: CreateEntry :one
INSERT INTO entries (
   account_id, amount
) VALUES (
  $1, $2
)
RETURNING *;

-- name: ListEntries :many
SELECT * FROM entries
ORDER BY account_id;


-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;

-- name: UpdateEntry :one
UPDATE entries
  set amount = $2
WHERE id = $1
RETURNING *;