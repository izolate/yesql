package yesql

import (
	"context"
	"testing"

	_ "github.com/lib/pq"
)

func TestExec(t *testing.T) {
	authors := []author{
		{Name: "George R.R. Martin"},
		{Name: "Lewis Carroll"},
		{Name: "J.R.R. Tolkien"},
		{Name: "Stephen King"},
		{Name: "Douglas Adams"},
		{Name: "Frank Herbert"},
		{Name: "George Orwell"},
	}

	for _, a := range authors {
		q := "INSERT INTO authors (name) VALUES (@Name);"
		res, err := db.Exec(q, a)
		if err != nil {
			t.Fatal(err)
		}
		ra, err := res.RowsAffected()
		if err != nil {
			t.Fatal(err)
		}
		if ra != int64(1) {
			t.Fatal("Incorrect rows affected")
		}
	}

	books := []book{
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

	for _, b := range books {
		q := "INSERT INTO books (title, author, genre) VALUES (@Title, @Author, @Genre);"
		res, err := db.ExecContext(context.TODO(), q, b)
		if err != nil {
			t.Fatal(err)
		}
		ra, err := res.RowsAffected()
		if err != nil {
			t.Fatal(err)
		}
		if ra != int64(1) {
			t.Fatal("Incorrect rows affected")
		}
	}
}
