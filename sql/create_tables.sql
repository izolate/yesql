-- Schema used to run unit tests
CREATE TABLE genres (
    id serial PRIMARY KEY,
    name text
);

INSERT INTO genres (name) VALUES ('Fantasy');
INSERT INTO genres (name) VALUES ('Horror');
INSERT INTO genres (name) VALUES ('Sci-Fi');
INSERT INTO genres (name) VALUES ('Self-Help');
INSERT INTO genres (name) VALUES ('Computers');
INSERT INTO genres (name) VALUES ('Comedy');

CREATE TABLE authors (
    id serial PRIMARY KEY,
    name text
);

CREATE TABLE books (
    id serial PRIMARY KEY,
    title text NOT NULL,
    author serial REFERENCES authors(id) NOT NULL,
    genre serial REFERENCES genres(id) NOT NULL
);