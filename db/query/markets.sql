-- name: CreateMarket :one
INSERT INTO markets (username, base_currency, quote_currency, min_order_amount, price_precision)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, username, base_currency, quote_currency, created_at;

-- name: GetMarketByID :one
SELECT id, username, base_currency, quote_currency, min_order_amount, price_precision, created_at
FROM markets
WHERE id = $1;

-- name: ListMarkets :many
SELECT id, base_currency, quote_currency, min_order_amount, price_precision, created_at
FROM markets
ORDER BY created_at DESC;

-- name: DeleteMarket :exec
DELETE FROM markets
WHERE id = $1;

-- name: GetMarketByCurrencies :one
SELECT id, username, base_currency, quote_currency, min_order_amount, price_precision, created_at
FROM markets
WHERE base_currency = $1 AND quote_currency = $2;

