CREATE TABLE IF NOT EXISTS connections
(
    id serial NOT NULL,
    port integer NOT NULL,
    encrypted_secret varchar NOT NULL,
    user_id integer NOT NULL references users(id) on delete no action,
    server_id integer NOT NULL references servers(id) on delete no action
);