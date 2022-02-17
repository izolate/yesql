package yesql

import (
	"context"
	"database/sql"

	"github.com/izolate/yesql/template"
)

type DB struct {
	DB  *sql.DB
	tpl template.Execer
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
		DB:  db,
		tpl: template.New(),
	}, nil
}

// ExecContext executes a query without returning any rows.
// The data object is a struct for any placeholder parameters in the query.
func (db *DB) ExecContext(ctx context.Context, query string, data interface{}) (sql.Result, error) {
	q, err := db.tpl.Exec(query, data)
	if err != nil {
		return nil, err
	}
	args := []interface{}{}
	return db.DB.ExecContext(ctx, q, args)
}

// Exec executes a query without returning any rows.
// The data object is a struct for any placeholder parameters in the query.
func (db *DB) Exec(query string, data interface{}) (sql.Result, error) {
	return db.ExecContext(context.Background(), query, data)
}

// QueryContext executes a query that returns rows, typically a SELECT.
// The data object is a struct for any placeholder parameters in the query.
func (db *DB) QueryContext(ctx context.Context, query string, data interface{}) (*sql.Rows, error) {
	q, err := db.tpl.Exec(query, data)
	if err != nil {
		return nil, err
	}
	args := []interface{}{}
	return db.DB.QueryContext(ctx, q, args)
}

// Query executes a query that returns rows, typically a SELECT.
// The data object is a struct for any placeholder parameters in the query.
func (db *DB) Query(ctx context.Context, query string, data interface{}) (*sql.Rows, error) {
	return db.QueryContext(context.Background(), query, data)
}

/*
import yesql

type User struct {
	Name string `db:"name"`
	Age  int    `db:"age"`
}

func getUserByName(ctx context.Context, u User) {
	db, err := yesql.Open("postgres", "...")
	stmt := `SELECT * FROM users WHERE name = :name {{if .Age}}AND age = :age{{end}}`
	rows, err := db.ExecContext(ctx, stmt, u)
}
*/
