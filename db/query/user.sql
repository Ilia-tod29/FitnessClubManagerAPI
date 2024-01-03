-- name: CreateUser :one
INSERT INTO users (
    email,
    hashed_password,
    suspended,
    role
) VALUES (
    $1, $2, $3, $4
) RETURNING id, email, suspended, role, created_at;

-- name: GetUser :one
SELECT id, email, suspended, role, created_at FROM users
WHERE id = $1 LIMIT 1;

-- name: ListAllUsers :many
SELECT id, email, suspended, role, created_at FROM users
ORDER BY id;

-- name: ListUsers :many
SELECT id, email, suspended, role, created_at FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
set suspended = $2
WHERE id = $1
RETURNING id, email, suspended, role, created_at;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING id, email, suspended, role, created_at;