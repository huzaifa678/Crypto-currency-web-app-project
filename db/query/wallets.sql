-- name: CreateWallet :one
INSERT INTO wallets (user_email, currency, balance)
VALUES ($1, $2, $3)
RETURNING id, user_email, currency, balance, locked_balance, created_at;

-- name: GetWalletByID :one
SELECT id, user_email, currency, balance, locked_balance, created_at
FROM wallets
WHERE id = $1;

-- name: UpdateWalletBalance :exec
UPDATE wallets
SET balance = $1, locked_balance = $2
WHERE id = $3;

-- name: DeleteWallet :exec
DELETE FROM wallets
WHERE id = $1;
