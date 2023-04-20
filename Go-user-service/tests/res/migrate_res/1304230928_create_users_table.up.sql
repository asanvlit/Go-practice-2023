CREATE TABLE account (
    id uuid PRIMARY KEY,
    email varchar(320) NOT NULL,
    password_hash varchar(64) NOT NULL,
    UNIQUE (email)
);
