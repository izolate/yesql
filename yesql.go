package yesql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/izolate/yesql/bindvar"
	"github.com/izolate/yesql/template"
)

type Execer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type Queryer interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type DB struct {
	DB   *sql.DB
	tpl  template.Execer
	bvar bindvar.Parser
}

// Open opens a database specified by its database driver name and a
// driver-specific data source name, usually consisting of at least a
// database name and connection information.
//
// Open may just validate its arguments without creating a connection
// to the database. To verify that the data source name is valid, call
// Ping.
//
// The returned DB is safe for concurrent use by multiple goroutines
// and maintains its own pool of idle connections. Thus, the Open
// function should be called just once. It is rarely necessary to
// close a DB.
func Open(driver, dsn string) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return &DB{
		DB:   db,
		tpl:  template.New(),
		bvar: bindvar.New(driver),
	}, nil
}

// ExecContext executes a query without returning any rows.
// The data object is a map/struct for any placeholder parameters in the query.
func (db *DB) ExecContext(ctx context.Context, query string, data interface{}) (sql.Result, error) {
	return execContext(db.DB, ctx, query, data, db.tpl, db.bvar)
}

// Exec executes a query without returning any rows.
// The data object is a map/struct for any placeholder parameters in the query.
func (db *DB) Exec(query string, data interface{}) (sql.Result, error) {
	return db.ExecContext(context.Background(), query, data)
}

// QueryContext executes a query that returns rows, typically a SELECT.
// The data object is a map/struct for any placeholder parameters in the query.
func (db *DB) QueryContext(ctx context.Context, query string, data interface{}) (*sql.Rows, error) {
	return queryContext(db.DB, ctx, query, data, db.tpl, db.bvar)
}

// Query executes a query that returns rows, typically a SELECT.
// The data object is a map/struct for any placeholder parameters in the query.
func (db *DB) Query(ctx context.Context, query string, data interface{}) (*sql.Rows, error) {
	return db.QueryContext(context.Background(), query, data)
}

// execContext executes a query without returning any rows.
// The data object is a map/struct for any placeholder parameters in the query.
func execContext(
	db Execer,
	ctx context.Context,
	query string,
	data interface{},
	tpl template.Execer,
	bvar bindvar.Parser,
) (sql.Result, error) {
	qt, err := tpl.Exec(query, data)
	if err != nil {
		return nil, fmt.Errorf("yesql: %s", err)
	}
	q, args, err := bvar.Parse(qt, data)
	if err != nil {
		return nil, fmt.Errorf("yesql: %s", err)
	}
	return db.ExecContext(ctx, q, args...)
}

// queryContext executes a query that returns rows, typically a SELECT.
// The data object is a map/struct for any placeholder parameters in the query.
func queryContext(
	db Queryer,
	ctx context.Context,
	query string,
	data interface{},
	tpl template.Execer,
	bvar bindvar.Parser,
) (*sql.Rows, error) {
	qt, err := tpl.Exec(query, data)
	if err != nil {
		return nil, fmt.Errorf("yesql: %s", err)
	}
	q, args, err := bvar.Parse(qt, data)
	if err != nil {
		return nil, fmt.Errorf("yesql: %s", err)
	}
	return db.QueryContext(ctx, q, args...)
}

/*
import yesql

type User struct {
	ID   string `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

type UserQuery struct {
	Name string
	Age  int
}

func getUserByName(ctx context.Context, q UserQuery) error {
	db, err := yesql.Open("postgres", "...")
	stmt := `SELECT * FROM users WHERE name = @Name {{if .Age}}AND age = @Age{{end}}`
	row, err := db.QueryContext(ctx, stmt, q)
	if err != nil {
		return err
	}
	var u User
	if err := yesql.Scan(row, &u); err != nil {
		return err
	}
}
*/
