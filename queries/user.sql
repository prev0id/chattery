-- name: CreateUser :exec
INSERT INTO users(login, password, username)
VALUES ($1, $2, $3);

-- name: UpdateUser :exec
UPDATE users
SET username = @new_username,
    login = @new_login,
    password = @new_password,
    avatar_id = @new_avatar_id,
    updated_at=now()
WHERE username = @old_username;

-- name: UserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: UserByLogin :one
SELECT * FROM users
WHERE login = $1;

-- name: DeleteUserByUsername :execrows
DELETE FROM users
WHERE username = $1;
