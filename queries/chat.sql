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

-- name: UserChatPreviewByType :many
SELECT
    c.id,
    c.type,
    c.name,
    c.created_at,
    c.updated_at,
    (m.id IS NOT NULL)::boolean           AS has_last_message,
    COALESCE(m.id, 0)                     AS last_message_id,
    COALESCE(m.user_id, 0)                AS last_message_user_id,
    COALESCE(m.text, '')                  AS last_message_text,
    COALESCE(m.created_at, TIMESTAMP 'epoch') AS last_message_created_at
FROM chats c
JOIN chat_participants cp
    ON cp.chat_id = c.id AND cp.user_id = $1
LEFT JOIN LATERAL (
    SELECT id, user_id, text, created_at
    FROM chat_messages
    WHERE chat_id = c.id
    ORDER BY created_at DESC, id DESC
    LIMIT 1
) m ON true
WHERE c.type = $2
ORDER BY m.created_at DESC NULLS LAST, c.updated_at DESC, c.id DESC;

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
