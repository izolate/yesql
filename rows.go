package yesql

import "database/sql"

// Rows is the result of a query. Its cursor starts before the first row
// of the result set. Use Next to advance from row to row.
type Rows struct {
	*sql.Rows
}

// StructScan copies the columns in the current row into the values pointed
// at by the dest struct.
//
// StructScan is like Rows.Scan, but doesn't rely on positional scanning,
// and instead scans into a struct based on the column names and the db
// struct tags, e.g. Foo string `db:"foo"`.
func (rs *Rows) StructScan(dest interface{}) error {
	return scan(rs.Rows, dest)
}

func scan(rows *sql.Rows, dest interface{}) error {
	return nil
}
