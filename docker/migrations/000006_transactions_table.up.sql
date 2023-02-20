CREATE TABLE transactions
(
    id           SERIAL UNIQUE                NOT NULL,
    from_account INT REFERENCES accounts (id),
    to_account   INT REFERENCES accounts (id) NOT NULL,
    amount       DECIMAL                      NOT NULL,
    status       VARCHAR(250)                 NOT NULL,
    date_created TIMESTAMP                    NOT NULL,
    date_updated TIMESTAMP
);
