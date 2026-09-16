-- +goose Up
create table chats (
    id serial primary key,
    created_at timestamp not null default now()
);

create table chat_users (
    id serial primary key,
    chat_id int not null references chats(id) on delete cascade,
    username text not null,
    unique (chat_id, username)
);

create table messages (
    id serial primary key,
    chat_id int references chats(id) on delete cascade,
    from_user text not null,
    text text not null,
    sent_at timestamp not null default now()
);

-- +goose Down
drop table messages;
drop table chat_users;
drop table chats;
