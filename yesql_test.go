package yesql

import (
	"testing"

	_ "github.com/lib/pq"
)

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

	books := []book{
		{Title: "A Storm Of Swords", Author: 1, Genre: 1},
		{Title: "Alice's Adventures In Wonderland", Author: 2, Genre: 1},
		{Title: "The Fellowship Of The Ring", Author: 3, Genre: 1},
	}

	for _, b := range books {
		q := "INSERT INTO books (title, author, genre) VALUES (@Title, @Author, @Genre);"
		res, err := db.Exec(q, b)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(res.RowsAffected())
	}

}
