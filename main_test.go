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

var (
	authors = []author{
		{Name: "George R.R. Martin"},
		{Name: "Lewis Carroll"},
		{Name: "J.R.R. Tolkien"},
		{Name: "Stephen King"},
		{Name: "Douglas Adams"},
		{Name: "Frank Herbert"},
		{Name: "George Orwell"},
	}

	books = []book{
		{Title: "A Storm Of Swords", Author: 1, Genre: 1},
		{Title: "Alice's Adventures In Wonderland", Author: 2, Genre: 1},
		{Title: "The Fellowship Of The Ring", Author: 3, Genre: 1},
		{Title: "Salem's Lot", Author: 4, Genre: 2},
		{Title: "It", Author: 4, Genre: 2},
		{Title: "The Shining", Author: 4, Genre: 2},
		{Title: "The Hitchhiker's Guide to the Galaxy", Author: 5, Genre: 3},
		{Title: "Dune", Author: 6, Genre: 3},
		{Title: "1984", Author: 7, Genre: 3},
	}
)

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

type assert struct {
	t *testing.T
}

func (a assert) NilErr(e error) {
	if e != nil {
		a.t.Fatalf("err is not nil: %s", e.Error())
	}
}

func (a assert) Truthy(b bool) {
	if !b {
		a.t.Fatal("is not true")
	}
}
func (a assert) StringEq(x, y string) {
	if x != y {
		a.t.Fatalf(`"%s" != "%s"`, y, x)
	}
}

func (a assert) IntEq(x, y int) {
	if x != y {
		a.t.Fatalf("expected %d, got %d", x, y)
	}
}
