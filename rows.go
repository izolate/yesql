package yesql

import (
	"database/sql"
	"fmt"
	"reflect"
)

// Rows is the result of a query. Its cursor starts before the first row
// of the result set. Use Next to advance from row to row.
type Rows struct {
	*sql.Rows
}

// ScanStruct copies the columns in the current row into the values pointed
// at by the dest struct.
//
// ScanStruct is like Rows.Scan, but doesn't rely on positional scanning,
// and instead scans into a struct based on the column names and the db
// struct tags, e.g. Foo string `db:"foo"`.
func (rs *Rows) ScanStruct(dest interface{}) error {
	return scan(rs.Rows, dest)
}

const structTagDB = "db"

func scan(rows *sql.Rows, dest interface{}) error {
	dv := reflect.ValueOf(dest)

	// Ensure destination is a pointer.
	k := dv.Kind()
	if k != reflect.Ptr {
		return fmt.Errorf("yesql: destination not a pointer: %s", k)
	}
	// Ensure destination elem is a struct.
	k = dv.Elem().Kind()
	if k != reflect.Struct {
		return fmt.Errorf("yesql: destination not a struct: %s", k)
	}

	// Identify the column names in the rows.
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	// Create a list of pointers to the fields in the destination struct
	// that correspond to the row columns, based on the db struct tag.
	dests := []interface{}{}
	for _, c := range cols {
		dfi, ok := fieldIndex(dv.Elem().Type(), structTagDB, c)
		if !ok {
			return fmt.Errorf("yesql: field not found in destination for column: %s", c)
		}
		df := dv.Elem().Field(dfi)
		dests = append(dests, df.Addr().Interface())
	}
	return rows.Scan(dests...)
}

// fieldIndex returns the index of a struct field that matches the db tag value.
func fieldIndex(t reflect.Type, key, value string) (int, bool) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Tag.Get(key) == value {
			return i, true
		}
	}
	return -1, false
}
