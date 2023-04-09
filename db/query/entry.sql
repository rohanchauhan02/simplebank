-- name: CreateEntries :one
INSERT INTO entries (
    account_id,
    account
) VALUES (
  $1, $2
) RETURNING *;


-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntry :many
SELECT * FROM entries
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateEntry :exec
UPDATE entries
SET account_id = $1
where id = $2;

-- name: DeleteEntry :exec
DELETE FROM entries
where id = $1;
