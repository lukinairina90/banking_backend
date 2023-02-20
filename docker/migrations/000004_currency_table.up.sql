CREATE TABLE currency
(
    id   SERIAL UNIQUE NOT NULL,
    name VARCHAR(50)   NOT NULL,
    code VARCHAR(3)    NOT NULL
);

INSERT INTO currency (name, code)
VALUES ('US Dollar', 'USD'),
       ('Ukrainian hryvnia', 'UAH'),
       ('Polish zloty', 'PLN');