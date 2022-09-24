CREATE TABLE IF NOT EXISTS servers
(
    id serial NOT NULL,
    ip_address varchar(13) NOT NULL,
    location varchar(2) NOT NULL,
    CHECK ( location <> lower(location) and char_length(location) = 2 )
)