-- name: CreateFee :one
INSERT INTO fees (market_id, maker_fee, taker_fee)
VALUES ($1, $2, $3)
RETURNING id, market_id, maker_fee, taker_fee, created_at;

-- name: GetFeeByMarketID :one
SELECT id, market_id, maker_fee, taker_fee, created_at
FROM fees
WHERE market_id = $1;

-- name: DeleteFee :exec
DELETE FROM fees
WHERE id = $1;
