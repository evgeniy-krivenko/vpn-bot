ALTER TABLE users ADD CONSTRAINT "users_id_unique" UNIQUE (id);
ALTER TABLE servers ADD CONSTRAINT "servers_id_unique" UNIQUE (id);

CREATE TABLE IF NOT EXISTS connections
(
    id serial NOT NULL,
    port integer NOT NULL,
    encrypted_secret varchar NOT NULL,
    user_id integer NOT NULL UNIQUE references users(id) on delete cascade,
    server_id integer NOT NULL UNIQUE references servers(id) on delete cascade
);