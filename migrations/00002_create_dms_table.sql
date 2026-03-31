-- +goose Up
-- +goose StatementBegin
CREATE TABLE dms (
    id              BIGSERIAL PRIMARY KEY,
    last_message_id BIGINT    NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE dm_participants (
    dm_id                BIGINT    NOT NULL,
    user_id              BIGINT    NOT NULL,
    last_read_message_id BIGINT    NOT NULL,
    created_at           TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE dm_messages (
    id          BIGSERIAL PRIMARY KEY,
    dm_id       BIGINT    NOT NULL,
    user_id     BIGINT    NOT NULL,
    text        TEXT      NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX dm_messages_dm_id_id_desc_idx
    ON dm_messages (dm_id, id DESC);

CREATE INDEX dm_messages_dm_id_created_id_desc_idx
    ON dm_messages (dm_id, created_at DESC, id DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dms;
DROP TABLE dm_participants;
DROP TABLE dm_messages;
-- +goose StatementEnd
