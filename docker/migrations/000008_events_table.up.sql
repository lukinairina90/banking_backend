CREATE TYPE event_type AS ENUM (
    'ACCOUNT_CREATED',
    'ACCOUNT_DELETED',
    'ACCOUNT_BLOCKED',
    'ACCOUNT_UNBLOCKED',
    'CARD_CREATED',
    'USER_BLOCKED',
    'USER_UNBLOCKED',
    'WITHDRAWAL',
    'DEPOSIT'
    );

CREATE TABLE event
(
    id       SERIAL UNIQUE NOT NULL,
    user_id  INT REFERENCES users (id),
    type     event_type    NOT NULL,
    metadata JSONB,
    time     TIMESTAMP
)