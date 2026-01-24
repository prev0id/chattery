-- name: CreateChat :one
INSERT INTO chats(type, name)
VALUES ($1, $2)
RETURNING id;

-- name: DeleteChat :exec
DELETE FROM chats
WHERE id=$1;

-- name: Chats :many
SELECT * FROM chats;

-- name: UserChats :many
SELECT * FROM chats
WHERE chats.id in (
    SELECT chat_id FROM chat_participants
    WHERE user_id=$1
);

-- name: AddParticipant :exec
INSERT INTO chat_participants(chat_id, user_id, role)
VALUES ($1, $2, $3);

-- name: DeleteParticipant :exec
DELETE FROM chat_participants
WHERE chat_id=$1 AND user_id=$2;

-- name: CreateMessage :one
INSERT INTO chat_messages(chat_id, user_id, text)
VALUES ($1, $2, $3)
RETURNING id;

-- name: FirstPageOfMessages :many
SELECT * FROM chat_messages
WHERE chat_id = $1
ORDER BY created_at DESC, id DESC
LIMIT $2;

-- name: NextPagesOfMessages :many
SELECT * FROM chat_messages
WHERE chat_id = $1 AND (created_at < $2 OR (created_at = $2 AND id < $3))
ORDER BY created_at DESC, id DESC
LIMIT $4;

-- name: ParticipantsForChat :many
SELECT * FROM chat_participants
WHERE chat_id = $1;
