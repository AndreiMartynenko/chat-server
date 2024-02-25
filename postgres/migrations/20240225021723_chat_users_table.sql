-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS chats (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    -- created_at timestamp not null default now(),
    -- updated_at timestamp
);

CREATE TABLE IF NOT EXISTS users_chats (
    chat_id BIGINT REFERENCES chats(id),
    user_id BIGINT NOT NULL CHECK (user_id > 0),
    PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE IF NOT EXISTS chats_messages (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT REFERENCES chats(id),
    user_id BIGINT NOT NULL CHECK (user_id > 0),
    text TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS chats_messages;
DROP TABLE IF EXISTS users_chats;
DROP TABLE IF EXISTS chats;
-- +goose StatementEnd

