-- name: CreateWallet :one
INSERT INTO wallets (user_email, currency, balance)
VALUES ($1, $2, $3)
RETURNING id, user_email, currency, balance, locked_balance, created_at;

-- name: GetWalletByUserEmailAndCurrency :one
SELECT user_email, currency, balance, locked_balance, created_at
FROM wallets
WHERE user_email = $1 AND currency = $2;

-- name: UpdateWalletBalance :exec
UPDATE wallets
SET balance = $1, locked_balance = $2
WHERE user_email = $3 AND currency = $4;

-- name: DeleteWallet :exec
DELETE FROM wallets
WHERE id = $1;
