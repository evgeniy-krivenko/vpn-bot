CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL      NOT NULL,
    user_id    INT,
    chat_id    INT,
    username   VARCHAR(60) NOT NULL,
    first_name VARCHAR(60),
    last_name  VARCHAR(60),
    chat_type  VARCHAR,
    created_at TIMESTAMP   NOT NULL DEFAULT NOW(),
    CHECK ( user_id IS NOT NULL OR chat_id IS NOT NULL)
);