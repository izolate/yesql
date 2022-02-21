package yesql

import (
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var db *DB

type category struct {
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

func TestMain(m *testing.M) {
	db = MustOpen("postgres", "host=localhost user=postgres password=postgres database=postgres sslmode=disable")
	if _, err := db.DB.Exec("DELETE FROM books"); err != nil {
		panic(err)
	}
	if _, err := db.DB.Exec("DELETE FROM authors"); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestExec(t *testing.T) {
	authors := []author{
		{Name: "George R.R. Martin"},
		{Name: "Lewis Carroll"},
		{Name: "J.R.R. Tolkien"},
	}

	for _, a := range authors {
		q := "INSERT INTO authors (name) VALUES (@Name);"
		res, err := db.Exec(q, a)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(res.RowsAffected())
	}
}
