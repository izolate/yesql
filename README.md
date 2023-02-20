# yesql

yesql is a lightweight Go library that adds convenience features to the standard library
package `database/sql`. It's specifically designed for writing raw SQL queries,
with a focus on providing greater efficiency and ease-of-use.

yesql leans on the standard library, keeping the same API, making it an ideal
drop-in replacement for `database/sql`.

It's a WIP and currently provides the following features:

* Named arguments (e.g. `SELECT * FROM foo WHERE id = @ID`)
* Templates for query building (cached and thread-safe)
* Statement logging
* Same API as `database/sql`

### The Elevator Pitch™

> yesql is like a child born from the union of `database/sql` and
> `text/templates`.

## Quick start

Start by opening a connection to a database and initializing a database driver:

```go
package foo

import (
    "github.com/izolate/yesql"
    _ "github.com/lib/pq"
)

func main() {
    db, err := yesql.Open("postgres", "host=localhost user=foo sslmode=disable")
    if err != nil {
        panic(err)
    }
}
```

### `Exec` / `ExecContext`

You can use the `Exec` function to execute a query without returning data. Named
parameters (`@Foo`) allow you to bind arguments to map (or struct) fields
without the risk of SQL injection.

```go
type Book struct {
    ID     string
    Title  string
    Author string
}

func InsertBook(b Book) error {
    q := `INSERT INTO users (id, title, author) VALUES (@ID, @Title, @Author)`
    _, err := db.Exec(q, b)
    return err
}
```

> :frog: NOTE: Named parameters only work with exported fields.

### `Query` / `QueryContext`

Use `Query` to execute a query and return data. Templates allow you to perform
complex logic without string concatenation or query building.

The templates alter the final SQL statement based on the input provided.

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

func SearchBooks(s BookSearch) ([]Book, error) {
    rows, err := db.Query(sqlSearchBooks, s)
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

Positional scanning is inflexible. Instead, use `db` struct tags and
`rows.ScanStruct()` to scan into a struct:

```go
type Book struct {
    ID     string `db:"id"`
    Title  string `db:"title"`
    Author string `db:"author"`
}

func SearchBooks(s BookSearch) ([]Book, error) {
    rows, err := db.Query(sqlSearchBooks, s)
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

## Configuration

You can easily configure yesql using a list of functional config setters, which
are accepted by all entry points. For instance, to disable statement logging,
you can pass in the `OptQuiet` setter as follows:

```go
db, err := yesql.Open("postgres", "host=localhost user=foo sslmode=disable", yesql.OptQuiet())
```

This will ensure that statement logs are not printed during execution.

## Feature checklist

- [x] Templated SQL statements
- [x] Named arguments (bindvars)
- [x] Statement logging
- [ ] Query tracing
- [x] Struct scanning
- [x] Unicode support
- [x] Postgres support
- [ ] Prepared statements
