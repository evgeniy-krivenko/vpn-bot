DROP TABLE IF EXISTS connections;

ALTER TABLE users DROP CONSTRAINT "users_id_unique";
ALTER TABLE servers DROP CONSTRAINT "servers_id_unique";