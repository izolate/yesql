package yesql

import (
	"database/sql"
	"errors"
)

// Row is the result of calling QueryRow to select a single row.
// It's a re-implementation of *sql.Row because the standard type is inaccessible.
type Row struct {
	// One of these two will be non-nil:
	err  error // deferred error for easy chaining
	rows *Rows
}

// scan is a generic scan that works for both Scan and StructScan.
//
// NOTE: It is a reimplementation of `sql.Row.Scan` from the standard lib package,
// because the `sql.Row` struct has private fields and is therefore inaccesible.
func (r *Row) scan(fn func(dest ...interface{}) error, dest ...interface{}) error {
	if r.err != nil {
		return r.err
	}

	// TODO(bradfitz): for now we need to defensively clone all
	// []byte that the driver returned (not permitting
	// *RawBytes in Rows.Scan), since we're about to close
	// the Rows in our defer, when we return from this function.
	// the contract with the driver.Next(...) interface is that it
	// can return slices into read-only temporary memory that's
	// only valid until the next Scan/Close. But the TODO is that
	// for a lot of drivers, this copy will be unnecessary. We
	// should provide an optional interface for drivers to
	// implement to say, "don't worry, the []bytes that I return
	// from Next will not be modified again." (for instance, if
	// they were obtained from the network anyway) But for now we
	// don't care.
	defer r.rows.Close()
	for _, dp := range dest {
		if _, ok := dp.(*sql.RawBytes); ok {
			return errors.New("yesql: RawBytes isn't allowed on Row.Scan")
		}
	}

	if !r.rows.Next() {
		if err := r.rows.Err(); err != nil {
			return err
		}
		return sql.ErrNoRows
	}

	if err := fn(dest...); err != nil {
		return err
	}
	// Make sure the query can be processed to completion with no errors.
	return r.rows.Close()

}

// ScanStruct copies the columns in the current row into the values pointed
// at by the dest struct.
//
// ScanStruct is like Rows.Scan, but doesn't rely on positional scanning,
// and instead scans into a struct based on the column names and the db
// struct tags, e.g. Foo string `db:"foo"`.
func (r *Row) ScanStruct(dest interface{}) error {
	fn := func(dest ...interface{}) error {
		if len(dest) == 1 {
			return r.rows.ScanStruct(dest[0])
		}
		return r.Scan(dest...)
	}
	return r.scan(fn, dest)
}

// Scan copies the columns from the matched row into the values
// pointed at by dest. See the documentation on Rows.Scan for details.
// If more than one row matches the query,
// Scan uses the first row and discards the rest. If no row matches
// the query, Scan returns ErrNoRows.
func (r *Row) Scan(dest ...interface{}) error {
	return r.scan(r.rows.Scan, dest...)
}

// Err provides a way for wrapping packages to check for
// query errors without calling Scan.
// Err returns the error, if any, that was encountered while running the query.
// If this error is not nil, this error will also be returned from Scan.
func (r *Row) Err() error {
	return r.err
}
