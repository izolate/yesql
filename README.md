# yesql

A tool to write SQL in Go more effectively. WIP.

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

Use `Exec` or `ExecContext` to execute a query without returning data. Named parameters (`:Foo`) allow you to bind arguments to map (or struct) fields:

```go
type User struct {
    ID    string
    Email string
}

func insertUser(c context.Context, u User) error {
    // ID and Email are bound to the User fields
    stmt := `INSERT INTO users (id, email) VALUES (:ID, :Email)`
    _, err := db.ExecContext(c, stmt, u)
    if err != nil {
        return err
    }
    return nil
}
```

## Feature checklist

- [x] Templated SQL statements
- [x] Named arguments (bindvars)
- [ ] Statement logging
- [ ] Query tracing
- [ ] Struct scanning
- [ ] Unicode support
- [x] Postgres support

## Todo

- [x] Exec/ExecContext
- [x] Query/QueryContext
- [ ] QueryRow/QueryRowContext