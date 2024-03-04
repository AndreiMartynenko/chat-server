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

-- +goose Down
-- +goose StatementBegin
DROP TABLE messages;
DROP TABLE chat_users;
DROP TABLE chats;
DROP TABLE users;
-- +goose StatementEnd
