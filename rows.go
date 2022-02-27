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

var (
// typeScanner = reflect.TypeOf((*sql.Scanner)(nil)).Elem()
// typeValuer  = reflect.TypeOf((*driver.Valuer)(nil)).Elem()
)

func scan(rows *sql.Rows, dest interface{}) error {
	v := reflect.ValueOf(dest)
	k := v.Kind()
	fmt.Println("Value:", v)
	fmt.Println("Kind:", k)
	// Ensure destination is a pointer.
	if k != reflect.Ptr {
		return fmt.Errorf("yesql: destination is not a pointer: %s", v.Kind())
	}
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	for _, c := range cols {
		// Find the field with the db tag that matches the column
		fmt.Println(c)
	}
	return nil
}
