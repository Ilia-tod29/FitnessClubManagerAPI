-- name: CreateGalleryItem :one
INSERT INTO gallery (
    image
) VALUES (
    $1
) RETURNING *;

-- name: GetGalleryItem :one
SELECT * FROM gallery
WHERE id = $1 LIMIT 1;

-- name: ListAllGalleryItems :many
SELECT * FROM gallery
ORDER BY id;

-- name: ListGalleryItems :many
SELECT * FROM gallery
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteGalleryItem :one
DELETE FROM gallery
WHERE id = $1
RETURNING *;