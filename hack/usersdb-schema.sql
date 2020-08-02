CREATE SCHEMA users;

CREATE TABLE users.users (
    email text PRIMARY KEY,
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    fullname text,
    password_hash text
);
