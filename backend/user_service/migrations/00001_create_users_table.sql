-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id          BIGSERIAL PRIMARY KEY,
    login       TEXT      NOT NULL UNIQUE,
    password    BYTEA     NOT NULL,
    username    TEXT      NOT NULL UNIQUE,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
