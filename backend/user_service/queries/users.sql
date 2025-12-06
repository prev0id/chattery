-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: InsertUser :one
INSERT INTO users(login, password, username)
VALUES ($1, $2, $3)
RETURNING id;

-- name: UpdateUser :exec
UPDATE users
SET login=$2,
    password=$3,
    username=$4,
    image_id=$5,
    updated_at=now()
WHERE id=$1;
