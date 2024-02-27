-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE chats (
    id SERIAL PRIMARY KEY
);

CREATE TABLE chat_users (
    chat_id INT REFERENCES chats(id),
    user_id INT REFERENCES users(id),
    PRIMARY KEY(chat_id, user_id)
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    chat_id INT REFERENCES chats(id),
    user_id INT REFERENCES users(id),
    text TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- CREATE TABLE IF NOT EXISTS chats (
--     id BIGSERIAL PRIMARY KEY,
--     title TEXT NOT NULL,
--     -- created_at timestamp not null default now(),
--     -- updated_at timestamp
-- );

-- CREATE TABLE IF NOT EXISTS users_chats (
--     chat_id BIGINT REFERENCES chats(id),
--     user_id BIGINT NOT NULL CHECK (user_id > 0),
--     PRIMARY KEY (chat_id, user_id)
-- );

-- CREATE TABLE IF NOT EXISTS chats_messages (
--     id BIGSERIAL PRIMARY KEY,
--     chat_id BIGINT REFERENCES chats(id),
--     user_id BIGINT NOT NULL CHECK (user_id > 0),
--     text TEXT NOT NULL,
--     timestamp TIMESTAMP NOT NULL
-- );
-- +goose StatementEnd

-- +goose Down
DROP TABLE messages;
DROP TABLE chat_users;
DROP TABLE chats;
DROP TABLE users;
-- +goose StatementEnd

