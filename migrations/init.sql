CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT                NOT NULL,
    coins         INT DEFAULT 1000
);