-- name: CreateServer :one
INSERT INTO servers(name)
VALUES ($1)
RETURNING id;

-- name: UpdateServer :exec
UPDATE servers
SET name=$2,
    updated_at=NOW()
WHERE id=$1;

-- name: GetUserServers :many
SELECT
    s.id AS id,
    s.name AS name,
    sp.created_at AS joined_at,
    sp.role AS role,
    t.id AS topic_id,
    t.name AS topic_name,
    t.type AS topic_type,
    t.created_at AS topic_created_at
FROM servers s
JOIN server_participants sp ON sp.server_id = s.id
LEFT JOIN topics t ON t.server_id = s.id
WHERE sp.user_id = $1
ORDER BY sp.created_at, s.id, t.created_at, t.id;

-- name: ListServers :many
SELECT
    s.id AS id,
    s.name AS name,
    t.id AS topic_id,
    t.name AS topic_name,
    t.type AS topic_type,
    t.created_at AS topic_created_at
FROM servers s
LEFT JOIN topics t ON t.server_id = s.id
ORDER BY s.id, t.created_at, t.id;

-- name: GetServer :many
SELECT
    s.id AS id,
    s.name AS name,
    t.id AS topic_id,
    t.name AS topic_name,
    t.type AS topic_type,
    t.created_at AS topic_created_at
FROM servers s
LEFT JOIN topics t ON t.server_id = s.id
WHERE s.id=$1
ORDER BY t.created_at, t.id;

-- name: CreateServerParticipant :exec
INSERT INTO server_participants (server_id, user_id, role)
VALUES ($1, $2, $3);

-- name: DeleteServerParticipant :exec
DELETE FROM server_participants
WHERE server_id = $1 AND user_id = $2;

-- name: GetServerParticipant :one
SELECT server_id, user_id, role, created_at, updated_at FROM server_participants
WHERE server_id = $1 AND user_id = $2;

-- name: GetServerParticipants :many
SELECT server_id, user_id, role, created_at, updated_at FROM server_participants
WHERE server_id = $1;

-- name: CreateTopic :one
INSERT INTO topics (server_id, name, type)
VALUES ($1, $2, $3)
RETURNING id;

-- name: CreateMessage :one
INSERT INTO topic_messages(topic_id, user_id, text)
VALUES ($1, $2, $3)
RETURNING id;

-- name: UpdateTopic :exec
UPDATE topics
SET
    name = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteTopic :exec
DELETE FROM topics
WHERE id = $1;

-- name: GetTopic :one
SELECT id, server_id, name, type, created_at, updated_at FROM topics
WHERE id=$1;

-- name: DeleteMessagesByServerID :exec
DELETE FROM topic_messages
WHERE topic_id IN (SELECT id FROM topics WHERE server_id = $1);

-- name: DeleteTopicsByServerID :exec
DELETE FROM topics
WHERE server_id = $1;

-- name: DeleteServerParticipantsByServerID :exec
DELETE FROM server_participants
WHERE server_id = $1;

-- name: DeleteServerByServerID :exec
DELETE FROM servers
WHERE id = $1;

-- name: FirstPageOfTopicMessages :many
SELECT id, topic_id, user_id, text, created_at, updated_at FROM topic_messages
WHERE topic_id = $1
ORDER BY created_at DESC, id DESC
LIMIT $2;

-- name: NextPagesOfTopicMessages :many
SELECT id, topic_id, user_id, text, created_at, updated_at FROM topic_messages
WHERE topic_id = $1 AND (created_at < $2 OR (created_at = $2 AND id < $3))
ORDER BY created_at DESC, id DESC
LIMIT $4;
