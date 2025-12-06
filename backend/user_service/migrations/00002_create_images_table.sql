-- +goose Up
-- +goose StatementBegin
CREATE TABLE images (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT    NOT NULL,
    image_id    TEXT      NOT NULL,
    file_type   TEXT      NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE images;
-- +goose StatementEnd
