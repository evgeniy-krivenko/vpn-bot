ALTER TABLE connections ADD CONSTRAINT "connections_user_id_key" UNIQUE (user_id);
ALTER TABLE connections ADD CONSTRAINT "connections_server_id_key" UNIQUE (server_id);