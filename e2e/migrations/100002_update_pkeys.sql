-- +goose Up
-- +goose StatementBegin
SELECT setval(
    pg_get_serial_sequence('users', 'id'),
    (SELECT max(id) FROM users)
);

SELECT setval(
    pg_get_serial_sequence('servers', 'id'),
    (SELECT max(id) FROM servers)
);

SELECT setval(
    pg_get_serial_sequence('topics', 'id'),
    (SELECT max(id) FROM topics)
);

SELECT setval(
    pg_get_serial_sequence('topic_messages', 'id'),
    (SELECT max(id) FROM topic_messages)
);

SELECT setval(
    pg_get_serial_sequence('dms', 'id'),
    (SELECT max(id) FROM dms)
);

SELECT setval(
    pg_get_serial_sequence('dm_messages', 'id'),
    (SELECT max(id) FROM dm_messages)
);
-- +goose StatementEnd
--
-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
