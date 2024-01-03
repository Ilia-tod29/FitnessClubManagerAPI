-- name: CreateSubscription :one
INSERT INTO subscriptions (
    user_id,
    start_date,
    end_date
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetSubscription :one
SELECT * FROM subscriptions
WHERE id = $1 LIMIT 1;

-- name: ListAllSubscriptions :many
SELECT * FROM subscriptions
ORDER BY id;

-- name: ListSubscriptions :many
SELECT * FROM subscriptions
ORDER BY id
LIMIT $1
OFFSET $2;
