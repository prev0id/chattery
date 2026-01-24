-- +goose Up
-- +goose StatementBegin
CREATE TABLE chats (
    id          BIGSERIAL PRIMARY KEY,
    type        TEXT      NOT NULL DEFAULT ''::TEXT,
    name        TEXT      NOT NULL DEFAULT ''::TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE chat_participants (
    chat_id     BIGINT    NOT NULL,
    user_id     BIGINT    NOT NULL,
    role        TEXT      NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE chat_messages (
    id          BIGSERIAL PRIMARY KEY,
    chat_id     BIGINT    NOT NULL,
    user_id     BIGINT    NOT NULL,
    text        TEXT      NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX chat_messages_chat_id_id_desc_idx
    ON chat_messages (chat_id, id DESC);

CREATE INDEX chat_messages_chat_id_created_id_desc_idx
    ON chat_messages (chat_id, created_at DESC, id DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE chats;
DROP TABLE chat_participants;
DROP TABLE chat_messages;
-- +goose StatementEnd
