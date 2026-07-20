-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (
    username,
    email,
    tenancy
) VALUES (
    sqlc.arg(username)::text,
    sqlc.arg(email)::text,
    sqlc.narg(tenancy)::text
) RETURNING *;


-- name: DeleteUser :exec
UPDATE users
SET
    deleted_at = NOW()
WHERE id = sqlc.arg(id)::uuid;


-- name: SearchUsers :many
SELECT * FROM users
WHERE 1=1
AND tenancy = sqlc.arg(tenancy)::text
AND (sqlc.narg(username_filters)::text[] IS NULL OR username = ANY(sqlc.narg(username_filters)::text[]))
AND (sqlc.narg(email_filters)::text[] IS NULL OR email = ANY(sqlc.narg(email_filters)::text[]))
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_val)::int
OFFSET sqlc.arg(offset_val)::int;