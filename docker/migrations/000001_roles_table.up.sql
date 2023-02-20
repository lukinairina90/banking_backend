CREATE TABLE roles
(
    id   INT UNIQUE  NOT NULL,
    name VARCHAR(30) NOT NULL
);

INSERT INTO roles (id, name)
VALUES (1, 'admin');
INSERT INTO roles (id, name)
VALUES (2, 'user');