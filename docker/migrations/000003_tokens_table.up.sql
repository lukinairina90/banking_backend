CREATE TABLE refresh_tokens
(
    id         SERIAL UNIQUE                               NOT NULL,
    user_id    INT REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    token      VARCHAR(255)                                NOT NULL UNIQUE,
    expires_at TIMESTAMP                                   NOT NULL
);