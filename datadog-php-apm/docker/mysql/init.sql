USE php-docker-db;

CREATE TABLE mytable (
    id INT PRIMARY KEY,
    name VARCHAR(255)
);

INSERT INTO mytable (id, name) VALUES (1, 'John Doe');