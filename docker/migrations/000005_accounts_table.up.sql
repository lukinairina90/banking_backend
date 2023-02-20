CREATE TABLE accounts
(
    id          SERIAL UNIQUE                NOT NULL,
    iban        VARCHAR(100)                 NOT NULL,
    user_id     INT REFERENCES users (id)    NOT NULL,
    currency_id INT REFERENCES currency (id) NOT NULL,
    blocked     BOOLEAN                      NOT NULL,
    amount      DECIMAL                      NOT NULL DEFAULT 0
);