package yesql

import (
	"os"
	"testing"
)

const (
	sqlUp = `
CREATE TABLE genres (
    id serial PRIMARY KEY,
    name text
);

INSERT INTO genres (name) VALUES ('Fantasy');
INSERT INTO genres (name) VALUES ('Horror');
INSERT INTO genres (name) VALUES ('Sci-Fi');

CREATE TABLE authors (
    id serial PRIMARY KEY,
    name text
);

CREATE TABLE books (
    id serial PRIMARY KEY,
    title text NOT NULL,
    author serial REFERENCES authors(id) NOT NULL,
    genre serial REFERENCES genres(id) NOT NULL
);`

	sqlDown = `
DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS authors;
DROP TABLE IF EXISTS genres;`
)

var db *DB

type genre struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type author struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type book struct {
	ID     int    `db:"id"`
	Title  string `db:"title"`
	Author int    `db:"author"`
	Genre  int    `db:"genre"`
}

func mustExec(q string) {
	if _, err := db.DB.Exec(q); err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	db = MustOpen("postgres", "host=localhost user=postgres password=postgres database=postgres sslmode=disable")
	mustExec(sqlDown)
	mustExec(sqlUp)
	os.Exit(m.Run())
}
