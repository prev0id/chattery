-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    username    TEXT      PRIMARY KEY,
    login       TEXT      NOT NULL UNIQUE,
    password    BYTEA     NOT NULL,
    avatar_id   TEXT      NOT NULL DEFAULT ''::text,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
