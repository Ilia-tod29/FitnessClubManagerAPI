-- name: CreateInventoryItem :one
INSERT INTO inventory (
    name,
    image
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetInventoryItem :one
SELECT * FROM inventory
WHERE id = $1 LIMIT 1;

-- name: ListAllInventoryItems :many
SELECT * FROM inventory
ORDER BY id;

-- name: ListInventoryItems :many
SELECT * FROM inventory
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateInventoryItem :one
UPDATE inventory
set name = $2,
    image = $3
WHERE id = $1
RETURNING *;

-- name: UpdateInventoryItemName :one
UPDATE inventory
set name = $2
WHERE id = $1
RETURNING *;

-- name: UpdateInventoryItemImage :one
UPDATE inventory
set image = $2
WHERE id = $1
RETURNING *;

-- name: DeleteInventoryItem :exec
DELETE FROM inventory
WHERE id = $1;