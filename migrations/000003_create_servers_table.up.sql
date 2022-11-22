CREATE TABLE IF NOT EXISTS servers
(
    id serial NOT NULL,
    ip_address varchar(13) NOT NULL,
    location varchar(6) NOT NULL
);