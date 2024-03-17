-- +goose Up
-- +goose StatementBegin

create table if not exists chats (
    id bigserial primary key
);

create table if not exists chats_users (
    chat_id bigserial references chats,
    user_id bigint not null,
    primary key (chat_id, user_id)
);

create table if not exists chats_messages (
    id bigserial primary key,
    chat_id bigint references chats,
    user_id bigint not null,
    body text not null,
    time timestamp with time zone not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists chats_messages;
drop table if exists chats_users;
drop table if exists chats;
-- +goose StatementEnd



