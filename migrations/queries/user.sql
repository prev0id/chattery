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
    updated_at=now()
WHERE id = $1;

-- name: UserByUsername :one
SELECT * FROM users
WHERE id = $1;

-- name: UserByLogin :one
SELECT * FROM users
WHERE login = $1;

-- name: DeleteUserByID :execrows
DELETE FROM users
WHERE id = $1;
