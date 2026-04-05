-- name: CreateDM :one
INSERT INTO dms(last_message_id)
VALUES (0)
RETURNING id;

-- name: SetLastMessageInDM :exec
UPDATE dms
SET last_message_id=$2,
    updated_at=now()
WHERE id = $1;

-- name: UserDMs :many
SELECT
    d.id AS dm_id,
    d.updated_at AS last_activity_at,
    p.last_read_message_id,
    lm.id AS last_message_id,
    lm.user_id AS last_message_sender_id,
    lm.text AS last_message_text,
    lm.created_at AS last_message_created_at,
    p2.user_id AS other_participant_id
FROM dm_participants p
JOIN dms d ON d.id = p.dm_id
LEFT JOIN dm_messages lm ON lm.id = d.last_message_id
JOIN dm_participants p2 ON p2.dm_id = d.id AND p2.user_id != $1
WHERE p.user_id = $1;

-- name: CreateDMParticipant :exec
INSERT INTO dm_participants(dm_id, user_id)
VALUES ($1, $2);

-- name: SetDMLastReadMessage :exec
UPDATE dm_participants
SET last_read_message_id=$3,
    updated_at=now()
WHERE dm_id=$1 AND user_id=$2;

-- name: CreateDMMessage :one
INSERT INTO dm_messages(dm_id, user_id, text)
VALUES ($1, $2, $3)
RETURNING id;

-- name: FirstPageOfDMMessages :many
SELECT * FROM dm_messages
WHERE dm_id = $1
ORDER BY created_at ASC, id ASC
LIMIT $2;

-- name: NextPagesOfDMMessages :many
SELECT * FROM dm_messages
WHERE dm_id = $1 AND (created_at < $2 OR (created_at = $2 AND id < $3))
ORDER BY created_at ASC, id ASC
LIMIT $4;

-- name: GetDMParticipant :one
SELECT * FROM dm_participants
WHERE dm_id = $1 AND user_id = $2;

-- name: GetDMBetweenUsers :one
SELECT d.id FROM dm_participants p1
JOIN dm_participants p2 ON p1.dm_id = p2.dm_id
JOIN dms d ON d.id = p1.dm_id
WHERE p1.user_id = $1 AND p2.user_id = $2
LIMIT 1;
