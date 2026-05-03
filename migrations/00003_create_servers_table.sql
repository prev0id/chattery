-- +goose Up
-- +goose StatementBegin

CREATE TABLE servers (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT      NOT NULL DEFAULT ''::TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE server_participants (
    server_id   BIGINT    NOT NULL,
    user_id     BIGINT    NOT NULL,
    role        TEXT      NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE topics (
    id          BIGSERIAL PRIMARY KEY,
    server_id   BIGINT    NOT NULL,
    name        TEXT      NOT NULL,
    type        TEXT      NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE topic_messages (
    id          BIGSERIAL PRIMARY KEY,
    topic_id    BIGINT    NOT NULL,
    user_id     BIGINT    NOT NULL,
    text        TEXT      NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE topic_messages;
DROP TABLE topics;
DROP TABLE server_participants;
DROP TABLE servers;
-- +goose StatementEnd
