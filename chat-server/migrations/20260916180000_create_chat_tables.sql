-- +goose Up
CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE chat_users (
    id SERIAL PRIMARY KEY,
    chat_id INT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    username TEXT NOT NULL,
    UNIQUE (chat_id, username)
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    chat_id INT REFERENCES chats(id) ON DELETE CASCADE,
    from_user TEXT NOT NULL,
    text TEXT NOT NULL,
    sent_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE messages;
DROP TABLE chat_users;
DROP TABLE chats;
