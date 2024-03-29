package yesql

import (
	"context"
	"database/sql"
)

// Tx is an in-progress database transaction.
//
// A transaction must end with a call to Commit or Rollback.
//
// After a call to Commit or Rollback, all operations on the
// transaction fail with ErrTxDone.
//
// The statements prepared for a transaction by calling
// the transaction's Prepare or Stmt methods are closed
// by the call to Commit or Rollback.
type Tx struct {
	*sql.Tx
	cfg *Config
}

// ExecContext executes a query that doesn't return rows.
// The data object is a map/struct for any placeholder parameters in the query.
func (tx *Tx) ExecContext(ctx context.Context, query string, data interface{}) (sql.Result, error) {
	return ExecContext(tx.Tx, ctx, query, data, tx.cfg)
}

// Exec executes a query without returning any rows.
// The data object is a map/struct for any placeholder parameters in the query.
func (tx *Tx) Exec(query string, data interface{}) (sql.Result, error) {
	return tx.ExecContext(context.Background(), query, data)
}

// QueryContext executes a query that returns rows, typically a SELECT.
// The data object is a map/struct for any placeholder parameters in the query.
func (tx *Tx) QueryContext(ctx context.Context, query string, data interface{}) (*Rows, error) {
	return QueryContext(tx.Tx, ctx, query, data, tx.cfg)
}

// Query executes a query that returns rows, typically a SELECT.
// The data object is a map/struct for any placeholder parameters in the query.
func (tx *Tx) Query(ctx context.Context, query string, data interface{}) (*Rows, error) {
	return tx.QueryContext(context.Background(), query, data)
}
