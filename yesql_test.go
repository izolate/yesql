package yesql

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

func TestExec(t *testing.T) {
	its := assert{t}
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
		its.NilErr(err)
		ra, err := res.RowsAffected()
		its.NilErr(err)
		its.IntEq(1, int(ra))
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
		q := `
		INSERT INTO books (
			title,
			author,
			genre
		) VALUES (
			@Title,
			@Author,
			@Genre
		);`
		res, err := db.ExecContext(context.TODO(), q, b)
		its.NilErr(err)
		ra, err := res.RowsAffected()
		its.NilErr(err)
		its.IntEq(1, int(ra))
	}
}

func TestQuery(t *testing.T) {
	t.Run("Templates", func(t *testing.T) {
		its := assert{t}
		type search struct {
			Title  string
			Author string
			Genre  string
		}
		tcs := []struct {
			search   search
			expected []string // expected titles
		}{
			{
				search: search{},
				expected: []string{
					"A Storm Of Swords",
					"Alice's Adventures In Wonderland",
					"The Fellowship Of The Ring",
					"Salem's Lot",
					"It",
					"The Shining",
					"The Hitchhiker's Guide to the Galaxy",
					"Dune",
					"1984",
				},
			},
			{
				search: search{Author: "Stephen King"},
				expected: []string{
					"Salem's Lot",
					"It",
					"The Shining",
				},
			},
			{
				search:   search{Author: "Stephen King", Title: "%salem%"},
				expected: []string{"Salem's Lot"},
			},
			{
				search: search{Genre: "Sci-Fi"},
				expected: []string{
					"The Hitchhiker's Guide to the Galaxy",
					"Dune",
					"1984",
				},
			},
			{
				search:   search{Genre: "Sci-Fi", Author: "Douglas Adams"},
				expected: []string{"The Hitchhiker's Guide to the Galaxy"},
			},
		}
		for _, tc := range tcs {
			q := `
			SELECT b.title
			FROM books b
			JOIN authors a ON a.id = b.author
			JOIN genres g ON g.id = b.genre
			WHERE true
			{{if .Title}}AND b.title ILIKE @Title{{end}}
			{{if .Author}}AND a.name = @Author{{end}}
			{{if .Genre}}AND g.name ILIKE @Genre{{end}}`
			rows, err := db.Query(q, tc.search)
			its.NilErr(err)
			result := []string{}
			for rows.Next() {
				var s string
				err := rows.Scan(&s)
				its.NilErr(err)
				result = append(result, s)
			}
			its.IntEq(len(tc.expected), len(result))
			for i, ex := range tc.expected {
				its.StringEq(ex, result[i])
			}
		}
	})

	t.Run("ScanStruct", func(t *testing.T) {
		its := assert{t}
		type entity struct {
			Book   string `db:"book"`
			Author string `db:"author"`
			Genre  string `db:"genre"`
		}
		q := `
		SELECT
			b.title AS book,
			a.name AS author,
			g.name AS genre
		FROM books b
		JOIN authors a ON a.id = b.author
		JOIN genres g ON g.id = b.genre`
		rows, err := db.QueryContext(context.TODO(), q, nil)
		its.NilErr(err)
		es := []entity{}
		for rows.Next() {
			var e entity
			fmt.Println("======= rows.Next ========")
			its.NilErr(rows.ScanStruct(&e))
			es = append(es, e)
		}
		its.IntEq(9, len(es))
		for _, e := range es {
			t.Log(e)
			// its.StringEq("foo", e.Author)
		}
	})
}
