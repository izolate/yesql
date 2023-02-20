package yesql

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

func TestExec(t *testing.T) {
	its := assert{t}

	for _, a := range authors {
		q := "INSERT INTO authors (name) VALUES (@Name);"
		res, err := db.Exec(q, a)
		its.NilErr(err)
		ra, err := res.RowsAffected()
		its.NilErr(err)
		its.IntEq(1, int(ra))
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
			ID     int    `db:"id"`
			Book   string `db:"book"`
			Author string `db:"author"`
			Genre  string `db:"genre"`
		}
		q := `
		SELECT
			b.id AS id,
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
			its.NilErr(rows.ScanStruct(&e))
			es = append(es, e)
		}
		its.IntEq(9, len(es))
		expected := []entity{
			{ID: 1, Book: "A Storm Of Swords", Author: "George R.R. Martin", Genre: "Fantasy"},
			{ID: 2, Book: "Alice's Adventures In Wonderland", Author: "Lewis Carroll", Genre: "Fantasy"},
			{ID: 3, Book: "The Fellowship Of The Ring", Author: "J.R.R. Tolkien", Genre: "Fantasy"},
			{ID: 4, Book: "Salem's Lot", Author: "Stephen King", Genre: "Horror"},
			{ID: 5, Book: "It", Author: "Stephen King", Genre: "Horror"},
			{ID: 6, Book: "The Shining", Author: "Stephen King", Genre: "Horror"},
			{ID: 7, Book: "The Hitchhiker's Guide to the Galaxy", Author: "Douglas Adams", Genre: "Sci-Fi"},
			{ID: 8, Book: "Dune", Author: "Frank Herbert", Genre: "Sci-Fi"},
			{ID: 9, Book: "1984", Author: "George Orwell", Genre: "Sci-Fi"},
		}
		for i, e := range es {
			t.Run(e.Book, func(t *testing.T) {
				its = assert{t}
				its.IntEq(expected[i].ID, e.ID)
				its.StringEq(expected[i].Book, e.Book)
				its.StringEq(expected[i].Author, e.Author)
				its.StringEq(expected[i].Genre, e.Genre)
			})
		}
	})
}

func TestQueryRow(t *testing.T) {
	t.Run("ScanStruct", func(t *testing.T) {
		its := assert{t}
		tcs := []struct {
			query    string
			data     interface{}
			expected string // expected titles
		}{
			{
				"SELECT * FROM books WHERE title ~* @title",
				map[string]interface{}{"title": "1984"},
				"1984",
			},
			{
				"SELECT * FROM books WHERE title ~* @Title",
				map[string]interface{}{"Title": "Alice"},
				"Alice's Adventures In Wonderland",
			},
			{
				"SELECT * FROM books WHERE title ~* @Title",
				struct{ Title string }{"fellowship of the ring"},
				"The Fellowship Of The Ring",
			},
			{
				"SELECT * FROM books WHERE author = @AuthorID",
				&struct{ AuthorID int }{6},
				"Dune",
			},
		}
		for _, tc := range tcs {
			var b book
			err := db.QueryRow(tc.query, tc.data).ScanStruct(&b)
			its.NilErr(err)
			its.StringEq(tc.expected, b.Title)
		}
	})

	t.Run("Scan", func(t *testing.T) {
		its := assert{t}
		tcs := []struct {
			query string
			data  interface{}
			id    int
			title string
		}{
			{
				"SELECT id, title FROM books WHERE title ~* @title",
				map[string]interface{}{"title": "1984"},
				9,
				"1984",
			},
			{
				"SELECT id, title FROM books WHERE title ~* @Title",
				map[string]interface{}{"Title": "Alice"},
				2,
				"Alice's Adventures In Wonderland",
			},
			{
				"SELECT id, title FROM books WHERE title ~* @Title",
				struct{ Title string }{"fellowship of the ring"},
				3,
				"The Fellowship Of The Ring",
			},
			{
				"SELECT id, title FROM books WHERE title ~* @Title",
				struct{ Title string }{"A"}, // match multiple rows
				1,
				"A Storm Of Swords",
			},
		}
		for _, tc := range tcs {
			var (
				id    int
				title string
			)
			err := db.QueryRow(tc.query, tc.data).Scan(&id, &title)
			its.NilErr(err)
			its.IntEq(tc.id, id)
			its.StringEq(tc.title, title)
		}
	})
}

func TestUnicode(t *testing.T) {
	testCases := []string{
		"😂😂😂😂😂",
		"The 👏  future 👏  of 👏  online 👏  conversation 👏  is 👏  a 👏  hand 👏  emoji 👏  clapping 👏  in 👏  your 👏  face 👏  forever.",
		"こんにちは世界",
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			its := assert{t}
			var s string
			err := db.QueryRow(fmt.Sprintf("SELECT '%s'", tc), nil).Scan(&s)
			its.NilErr(err)
			its.StringEq(tc, s)
		})
	}
}
