-- name: CreateAccount :execresult
INSERT INTO accounts (
  owner,
  balance,
  currency
) VALUES (
  ?, ?, ?
);

-- name: GetAccouont :one
SELECT * FROM accounts
WHERE id = ? LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT ?
OFFSET ?;

-- name: UpdateAccount :exec
UPDATE accounts
SET balance = ?
WHERE id = ?;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = ?;