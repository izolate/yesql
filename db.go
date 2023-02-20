package yesql

import (
	"context"
	"database/sql"
)

type DB struct {
	*sql.DB
	cfg *Config
}

// ExecContext executes a query without returning any rows, e.g. an INSERT.
// The data object is a map/struct for any placeholder parameters in the query.
func (db *DB) ExecContext(ctx context.Context, query string, data interface{}) (sql.Result, error) {
	return ExecContext(db.DB, ctx, query, data, db.cfg)
}

// Exec executes a query without returning any rows, e.g. an INSERT.
// The data object is a map/struct for any placeholder parameters in the query.
func (db *DB) Exec(query string, data interface{}) (sql.Result, error) {
	return db.ExecContext(context.Background(), query, data)
}

// QueryContext executes a query that returns rows, typically a SELECT.
// The data object is a map/struct for any placeholder parameters in the query.
func (db *DB) QueryContext(ctx context.Context, query string, data interface{}) (*Rows, error) {
	return QueryContext(db.DB, ctx, query, data, db.cfg)
}

// Query executes a query that returns rows, typically a SELECT.
// The data object is a map/struct for any placeholder parameters in the query.
func (db *DB) Query(query string, data interface{}) (*Rows, error) {
	return db.QueryContext(context.Background(), query, data)
}

// QueryRowContext executes a query that is expected to return at most one row.
// QueryRowContext always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
func (db *DB) QueryRowContext(ctx context.Context, query string, data interface{}) *Row {
	return QueryRowContext(db.DB, ctx, query, data, db.cfg)
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
//
// QueryRow uses context.Background internally; to specify the context, use
// QueryRowContext.
func (db *DB) QueryRow(query string, data interface{}) *Row {
	return db.QueryRowContext(context.Background(), query, data)
}

// BeginTx starts a transaction.
//
// The provided context is used until the transaction is committed or rolled back.
// If the context is canceled, the sql package will roll back
// the transaction. Tx.Commit will return an error if the context provided to
// BeginTx is canceled.
//
// The provided TxOptions is optional and may be nil if defaults should be used.
// If a non-default isolation level is used that the driver doesn't support,
// an error will be returned.
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{tx, db.cfg}, nil
}

// Begin starts a transaction. The default isolation level is dependent on
// the driver.
func (db *DB) Begin() (*Tx, error) {
	return db.BeginTx(context.Background(), nil)
}
