# yesql

A tool to write raw SQL in Go more effectively. WIP.

Consider it a thin client over, or a drop-in replacement for, the standard library package `database/sql`.

### Features
* Same API as the standard library
* Named arguments support
* Templates for query building
* Query logs (TODO)

## Quick start

Open a connection to a database:

```go
package foo

import "github.com/izolate/yesql"

func main() {
    db, err := yesql.Open("postgres", "host=localhost user=foo sslmode=disable")
    if err != nil {
        panic(err)
    }
}
```

### `Exec`

Use `Exec` or `ExecContext` to execute a query without returning data. Named parameters (`@Foo`) allow you to bind arguments to map (or struct) fields without the risk of SQLi:

```go
type Book struct {
    ID     string
    Title  string
    Author string
}

func InsertBook(ctx context.Context, b Book) error {
    q := `INSERT INTO users (id, title, author) VALUES (@ID, @Title, @Author)`
    _, err := db.ExecContext(ctx, q, b)
    return err
}
```

### `Query`

Use `Query` or `QueryContext` to execute a query to return data. Templates offer you the chance to perform complex logic without string concatenation or query building:

```go
type BookSearch struct {
    Author string    
    Title  string
    Genre  string
}

const sqlSearchBooks = `
SELECT * FROM books
WHERE author = @Author
{{if .Title}}AND title ILIKE @Title{{end}}
{{if .Genre}}AND genre = @Genre{{end}}
`

func SearchBooks(ctx context.Context, s BookSearch) ([]Book, error) {
    rows, err := db.QueryContext(ctx, sqlSearchBooks, s)
    if err != nil {
        return nil, err
    }
    books := []Book{}
    for rows.Next() {
        var b Book
        if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Genre); err != nil {
            return nil, err
        }
        books = append(books, b)
    }
    return books, nil
}
```

In fact, positional scanning is inflexible. Let's scan into a struct instead using `db` struct tags and `rows.ScanStruct()`:

```go
type Book struct {
    ID     string `db:"id"`
    Title  string `db:"title"`
    Author string `db:"author"`
}

func SearchBooks(ctx context.Context, s BookSearch) ([]Book, error) {
    rows, err := db.QueryContext(ctx, sqlSearchBooks, s)
    if err != nil {
        return nil, err
    }
    books := []Book{}
    for rows.Next() {
        var b Book
        if err := rows.ScanStruct(&b); err != nil {
            return nil, err
        }
        books = append(books, b)
    }
    return books, nil
}
```

## Feature checklist

- [x] Templated SQL statements
- [x] Named arguments (bindvars)
- [ ] Statement logging
- [ ] Query tracing
- [x] Struct scanning
- [ ] Unicode support
- [x] Postgres support

## TODO

- [x] Exec/ExecContext
- [x] Query/QueryContext
- [x] QueryRow/QueryRowContext
- [ ] ScanSlice/ScanMap
