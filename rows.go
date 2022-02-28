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

const tagKeyDB = "db"

var (
// typeScanner = reflect.TypeOf((*sql.Scanner)(nil)).Elem()
// typeValuer  = reflect.TypeOf((*driver.Valuer)(nil)).Elem()
)

func scan(rows *sql.Rows, dest interface{}) error {
	dv := reflect.ValueOf(dest)

	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	dests := []interface{}{}
	for _, c := range cols {
		dfi, ok := fieldIndex(dv.Elem().Type(), tagKeyDB, c)
		if !ok {
			return fmt.Errorf("yesql: field not found in destination for column: %s", c)
		}
		fmt.Println("DFI:", dfi)
		df := dv.Elem().Field(dfi)
		dests = append(dests, df.Addr().Interface())
	}

	fmt.Println(dests)
	// Loop over dest fields - START
	// Loop over dest fields - END
	if err := rows.Scan(dests...); err != nil {
		return err
	}
	fmt.Println("RESULT", dest)
	return nil
}

func scanWORKS(rows *sql.Rows, dest interface{}) error {
	type A struct {
		A string
	}
	a := A{}
	av := reflect.ValueOf(&a).Elem().Field(0).Addr()
	// av := reflect.ValueOf(&a.A)
	var b, c string
	if err := rows.Scan(av.Interface(), &b, &c); err != nil {
		return err
	}
	fmt.Println("RESULT", a, b, c)
	return nil
}

func scan3(rows *sql.Rows, dest interface{}) error {
	var a, b, c string
	av := reflect.ValueOf(&a)
	if err := rows.Scan(av.Interface(), &b, &c); err != nil {
		return err
	}
	fmt.Println("RESULT", a, b, c)
	return nil
}

func scan2(rows *sql.Rows, dest interface{}) error {
	t := reflect.TypeOf(dest)
	fmt.Println("Type:", t)

	// Ensure destination is a pointer.
	k := t.Kind()
	if k != reflect.Ptr {
		return fmt.Errorf("yesql: destination not a pointer: %s", k)
	}

	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	dests := []interface{}{}
	for _, c := range cols {
		// Find the field with the db tag that matches the column name.
		f, ok := field(t.Elem(), tagKeyDB, c)
		if !ok {
			return fmt.Errorf("yesql: field not found in destination for column: %s", c)
		}
		fmt.Println("FIELD FOUND:", f)
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

// field finds a field in a struct with a struct tag key that that equals a value.
// e.g. Name string `db:"name"` => field(t, "db", "name") => Name
func field(t reflect.Type, key, value string) (reflect.StructField, bool) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Tag.Get(key) == value {
			return f, true
		}
	}
	return reflect.StructField{}, false
}
