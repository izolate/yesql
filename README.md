# yesql

**Write SQL. Skip the plumbing.**

yesql is a lightweight layer over Go's `database/sql` for writing raw SQL with
named parameters, cached templates, structured logging, and struct scanning.

> `database/sql` meets `text/template`—with less query plumbing.

## Quick start

Open a database connection with your usual driver:

```go
package main

import (
    "github.com/izolate/yesql"
    _ "github.com/lib/pq"
)

func main() {
    db, err := yesql.Open("postgres", "host=localhost user=foo sslmode=disable")
    if err != nil {
        panic(err)
    }
    defer db.Close()
}
```

Then keep the query in SQL. Named parameters bind values safely, templates add
optional clauses, and `ScanStruct` maps result columns to `db` tags:

```go
type Book struct {
    ID     string `db:"id"`
    Title  string `db:"title"`
    Author string `db:"author"`
    Genre  string `db:"genre"`
}

type BookSearch struct {
    Author string
    Title  string
    Genre  string
}

const searchBooksSQL = `
SELECT id, title, author, genre
FROM books
WHERE author = @Author
{{if .Title}}AND title ILIKE @Title{{end}}
{{if .Genre}}AND genre = @Genre{{end}}
`

func SearchBooks(db *yesql.DB, search BookSearch) ([]Book, error) {
    rows, err := db.Query(searchBooksSQL, search)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var books []Book
    for rows.Next() {
        var book Book
        if err := rows.ScanStruct(&book); err != nil {
            return nil, err
        }
        books = append(books, book)
    }
    return books, rows.Err()
}
```

Named parameters can bind from maps or exported struct fields.

## Configuration

yesql accepts functional options at setup. For example, `OptQuiet` disables
statement logging:

```go
db, err := yesql.Open(
    "postgres",
    "host=localhost user=foo sslmode=disable",
    yesql.OptQuiet(),
)
```

## Status

yesql is a work in progress.

Available today:

- Templated SQL statements
- Named parameters
- Structured statement logging
- Struct scanning
- Unicode support
- PostgreSQL support

Planned:

- Query tracing
- Prepared statements
