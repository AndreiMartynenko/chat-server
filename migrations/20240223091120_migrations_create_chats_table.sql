-- +goose Up
-- +goose StatementBegin
CREATE TABLE chats (
    id uuid primary key default gen_random_uuid(),
    usernames text[]
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS chats;
-- +goose StatementEnd
