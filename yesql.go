package yesql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/izolate/yesql/bindvar"
	"github.com/izolate/yesql/template"
)

type Execer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type Queryer interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// ExecerQueryer is a union interface comprising Execer and Queryer.
type ExecerQueryer interface {
	Execer
	Queryer
}

// New instantiates yesql with an existing database connection.
func New(db *sql.DB) (*DB, error) {
	drivers := sql.Drivers()
	if len(drivers) == 0 {
		return nil, errors.New("yesql: no sql driver found")
	}
	return &DB{
		DB:   db,
		tpl:  template.New(),
		bvar: bindvar.New(drivers[0]),
	}, nil
}

// Open opens a database specified by its database driver name and a
// driver-specific data source name, usually consisting of at least a
// database name and connection information.
//
// Most users will open a database via a driver-specific connection
// helper function that returns a *DB. No database drivers are included
// in the Go standard library. See https://golang.org/s/sqldrivers for
// a list of third-party drivers.
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

// MustOpen opens a database specified by its database driver name and a
// driver-specific data source name, and panics on any errors.
func MustOpen(driver, dsn string) *DB {
	db, err := Open(driver, dsn)
	if err != nil {
		panic(err)
	}
	return db
}

// ExecContext executes a query without returning any rows, e.g. an INSERT.
// The data object is a map/struct for any placeholder parameters in the query.
func ExecContext(
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

// QueryContext executes a query that returns rows, typically a SELECT.
// The data object is a map/struct for any placeholder parameters in the query.
func QueryContext(
	db Queryer,
	ctx context.Context,
	query string,
	data interface{},
	tpl template.Execer,
	bvar bindvar.Parser,
) (*Rows, error) {
	qt, err := tpl.Exec(query, data)
	if err != nil {
		return nil, fmt.Errorf("yesql: %s", err)
	}
	q, args, err := bvar.Parse(qt, data)
	if err != nil {
		return nil, fmt.Errorf("yesql: %s", err)
	}
	rows, err := db.QueryContext(ctx, q, args...)
	return &Rows{rows}, err
}

// QueryRowContext executes a query that is expected to return at most one row.
// QueryRowContext always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
func QueryRowContext(
	db Queryer,
	ctx context.Context,
	query string,
	data interface{},
	tpl template.Execer,
	bvar bindvar.Parser,
) *Row {
	rows, err := QueryContext(db, ctx, query, data, tpl, bvar)
	return &Row{rows: rows, err: err}
}
