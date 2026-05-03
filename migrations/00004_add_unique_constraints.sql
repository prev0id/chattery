-- +goose Up
-- +goose StatementBegin
ALTER TABLE dm_participants
    ADD CONSTRAINT dm_participants_dm_id_user_id_key UNIQUE (dm_id, user_id);

ALTER TABLE server_participants
    ADD CONSTRAINT server_participants_server_id_user_id_key UNIQUE (server_id, user_id);

ALTER TABLE topics
    ADD CONSTRAINT topics_server_id_name_type_key UNIQUE (server_id, name, type);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE topics
    DROP CONSTRAINT topics_server_id_name_type_key;

ALTER TABLE server_participants
    DROP CONSTRAINT server_participants_server_id_user_id_key;

ALTER TABLE dm_participants
    DROP CONSTRAINT dm_participants_dm_id_user_id_key;
-- +goose StatementEnd
