-- name: CreateUser :one
INSERT INTO users (
    email,
    password_hash,
    role
) VALUES (
    $1, $2, $3
)
RETURNING id, email, role, is_active, created_at, updated_at;


-- name: GetUserByID :one
SELECT id, email, password_hash, role, is_active,  created_at, updated_at
FROM users
WHERE id = $1
LIMIT 1;


-- name: GetUserByEmail :one
SELECT id, email, password_hash, role, is_active,  created_at, updated_at
FROM users
WHERE email = $1
LIMIT 1;


-- name: ListUsers :many
SELECT id, email, role, is_active, created_at, updated_at
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;


-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;