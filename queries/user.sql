-- name: CreateUser :one
INSERT INTO users(login, password, username)
VALUES ($1, $2, $3)
RETURNING id;

-- name: UpdateUser :exec
UPDATE users
SET username=$2,
    login=$3,
    password=$4,
    avatar_id=$5,
    updated_at=NOW()
WHERE id = $1;

-- name: UserByID :one
SELECT id, username, login, password, avatar_id, created_at, updated_at FROM users
WHERE id = $1;

-- name: UserByLogin :one
SELECT id, username, login, password, avatar_id, created_at, updated_at FROM users
WHERE login = $1;

-- name: DeleteUserByID :execrows
DELETE FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT id, username, login, password, avatar_id, created_at, updated_at FROM users;
