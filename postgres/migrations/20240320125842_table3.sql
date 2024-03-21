-- +goose Up
-- +goose StatementBegin

create table if not exists chats (
                                     id serial primary key
);

create table if not exists chats_users (
                                           chat_id serial references chats,
                                           user_id serial not null,
                                           primary key (chat_id, user_id)
);

create table if not exists chats_messages (
                                              id serial primary key,
                                              chat_id serial references chats,
                                              user_id serial not null,
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



