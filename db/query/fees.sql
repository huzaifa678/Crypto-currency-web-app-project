-- name: CreateFee :one
INSERT INTO fees (username, market_id, maker_fee, taker_fee)
VALUES ($1, $2, $3, $4)
RETURNING id, market_id, maker_fee, taker_fee, created_at;

-- name: GetFeeByMarketID :one
SELECT id, username, market_id, maker_fee, taker_fee, created_at
FROM fees
WHERE market_id = $1;

-- name: DeleteFee :exec
DELETE FROM fees
WHERE id = $1;
