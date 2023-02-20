CREATE TABLE users
(
    id            SERIAL UNIQUE             NOT NULL,
    name          VARCHAR(100)              NOT NULL,
    surname       VARCHAR(100)              NOT NULL,
    email         VARCHAR(100)              NOT NULL UNIQUE,
    password      VARCHAR(255)              NOT NULL,
    role_id       INT REFERENCES roles (id) NOT NULL,
    blocked       BOOLEAN                   NOT NULL,
    registered_at TIMESTAMP                 NOT NULL
);

INSERT INTO users (name, surname, email, password, role_id, blocked, registered_at)
VALUES ('admin', 'admin', 'example@gmail.com', '73616c74d033e22ae348aeb5660fc2140aec35850c4da997',
        (select id from roles where name = 'admin'), false, now() AT TIME ZONE 'utc')
;