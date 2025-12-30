-- name: CreateUser :exec
INSERT INTO users(login, password, username)
VALUES ($1, $2, $3);

-- name: UpdateUser :exec
UPDATE users
SET login=$2,
    password=$3,
    username=$4,
    avatar_id=$5,
    updated_at=now()
WHERE username=$1;

-- name: UserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: UserByLogin :one
SELECT * FROM users
WHERE login = $1;
