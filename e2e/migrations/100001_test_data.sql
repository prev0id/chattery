-- +goose Up
-- +goose StatementBegin
INSERT INTO users (id, username, login, password, avatar_id) VALUES
(1, 'alex', 'alex@example.com', '$2a$10$djnwZqk0Jh2w4Ej.Mlk2desVJWcE1ZaW3oqnJG4oi1.HE8ij5JuR.', ''),
(2, 'bob', 'bob@example.com', '$2a$10$ScW7NGj/4vgLQ8t4cQ2EVeJxKMLL0uju5ylvd6xMMW/bZ4gclJzJG', ''),
(3, 'charlie', 'charlie@example.com', '$2a$10$LNzRuFpStCaeASUIDbAf5./EvXBLQhlgarS5Z/O3WK7rplcLXetc2', '');

INSERT INTO servers (id, name) VALUES
(1, 'Empty Server'),
(2, 'Work Space'),
(3, 'Gaming Hub'),
(4, 'Music Club');

INSERT INTO server_participants (server_id, user_id, role) VALUES
(2, 1, 'owner'),
(2, 2, 'member'),
(3, 1, 'member'),
(3, 2, 'owner'),
(4, 1, 'owner'),
(4, 2, 'member');

INSERT INTO topics (id, server_id, name, type) VALUES
(1, 2, 'general', 'text'),
(2, 2, 'random', 'voice'),
(3, 2, 'announcements', 'text'),
(4, 3, 'minecraft', 'text'),
(5, 3, 'cs2', 'voice'),
(6, 3, 'Dota2', 'text'),
(7, 4, 'rock', 'text'),
(8, 4, 'jazz', 'voice'),
(9, 4, 'classical', 'text');

INSERT INTO topic_messages (topic_id, user_id, text) VALUES
(1, 1, 'Welcome to Work Space!'),
(1, 2, 'Thanks, happy to be here'),
(1, 1, 'Lets get some work done'),
(2, 2, 'Anyone for a voice chat?'),
(3, 1, 'Please read the rules'),
(4, 2, 'Who wants to play Minecraft?'),
(4, 1, 'Count me in!'),
(5, 1, 'Looking for teammates'),
(5, 2, 'I am in'),
(6, 2, 'Anyone up for Dota?'),
(7, 1, 'What is your favorite rock band?'),
(7, 2, 'Pink Floyd'),
(8, 2, 'Jazz session tonight'),
(9, 1, 'Beethoven is the best');

INSERT INTO dms (id, last_message_id) VALUES
(1, 5);

INSERT INTO dm_participants (dm_id, user_id, last_read_message_id) VALUES
(1, 1, 5),
(1, 2, 5);

INSERT INTO dm_messages (dm_id, user_id, text) VALUES
(1, 1, 'Hey bob, how are you?'),
(1, 2, 'I am doing great!'),
(1, 1, 'Wanna play some games later?'),
(1, 2, 'Sure, what do you have in mind?'),
(1, 1, 'Maybe some Minecraft?');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM topic_messages;
DELETE FROM topics;
DELETE FROM server_participants;
DELETE FROM servers;
DELETE FROM dm_messages;
DELETE FROM dm_participants;
DELETE FROM dms;
DELETE FROM users;
-- +goose StatementEnd
