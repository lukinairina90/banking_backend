CREATE TABLE cards
(
    id              SERIAL UNIQUE NOT NULL,
    account_id      INT REFERENCES accounts (id) ON DELETE CASCADE,
    card_number     VARCHAR(16)   NOT NULL,
    cardholder_name VARCHAR(250)  NOT NULL,
    expiration_date TIMESTAMP     NOT NULL,
    cvv_code        VARCHAR(3)    NOT NULL
);