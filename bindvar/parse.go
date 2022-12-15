package bindvar

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
)

// Named argument prefix syntax used by the std lib.
// https://pkg.go.dev/database/sql#Named
const naPrefix = "@"

type Parser interface {
	// Parse parses all named parameters in a SQL statement, and returns
	// a statement with the params converted to bindvars appropriate for
	// the engine, e.g. :foo, :bar => $1, $2 (postgres).
	//
	// Additionally, the positional args are returned in order.
	Parse(query string, data interface{}) (q string, args []interface{}, err error)
}

// New creates a new parser.
func New(driver string) Parser {
	return &parser{driver}
}

type parser struct {
	driver string
}

func (p parser) Parse(query string, data interface{}) (string, []interface{}, error) {
	// Convert to rune to handle unicode strings
	qt := []rune(query)

	// Parse named args
	q, nvs := parse(p.driver, qt)
	args := []interface{}{}
	for _, nv := range nvs {
		// Get the named arg values from data
		v := value(data, nv.Name)
		args = append(args, v)
	}

	return string(q), args, nil
}

// reArgTerm is the terminating character of a named arg.
var reArgTerm = regexp.MustCompile(`[[:space:]]|;|\)|,`)

// parse parses the named args out of a query and returns a string with
// the correct arg syntax for the driver, and a list of arg names.
func parse(driverName string, query []rune) (s []rune, args []driver.NamedValue) {
	var (
		a      int  // Pointer used to seek through the string
		op     int  // The ordinal position of the captured arg
		ignore bool // Used to ignore false positives
	)
	for a < len(query) {
		ra := query[a] // the rune at position a

		// Ignore characters inside sql string literals.
		if string(ra) == "'" {
			ignore = !ignore
		}

		if !ignore && string(ra) == naPrefix {
			// We've found an argument! Create second pointer to find end of argument.
			b := a

			// Find the first terminating character to infer the end of the arg.
			for b < len(query) {
				rb := query[b]
				if reArgTerm.MatchString(string(rb)) {
					break
				}
				b++
			}

			op++ // Increment the arg's ordinal position.

			// Get the name of the arg, ignoring the prefix (@).
			a1 := a + 1
			n := string(query[a1:b])

			// Add the named arg to the list of all found args.
			nv := driver.NamedValue{
				Ordinal: op,
				Name:    n,
			}
			args = append(args, nv)

			// Convert the named arg to the correct syntax for the driver.
			arg := []rune(argfmt(driverName, nv))
			s = append(s, arg...)

			a = b // Skip to the end of the arg
			continue
		}

		s = append(s, query[a])
		a++
	}
	return s, args
}

// value gets the value for field (name) in the data object.
func value(data interface{}, name string) interface{} {
	if m, ok := data.(map[string]interface{}); ok {
		if v, ok := m[name]; ok {
			return v
		}
		return nil
	}

	// If data is not a simple map, use reflection to get the value.
	v := reflect.Indirect(reflect.ValueOf(data))
	switch {
	case v.Kind() == reflect.Struct: // Struct
		if f := v.FieldByName(name); f.IsValid() {
			return f.Interface()
		}
	case v.Elem().Kind() == reflect.Struct: // Pointer struct
		el := v.Elem()
		if f := el.FieldByName(name); f.IsValid() {
			return f.Interface()
		}
	case v.Elem().Kind() == reflect.Map: // Map pointer
		if val := v.Elem().MapIndex(reflect.ValueOf(name)); !val.IsZero() {
			return val.Interface()
		}
	}
	return nil
}

// argfmt converts a named arg to the correct syntax for the driver.
// e.g. @Foo => $1 (postgres)
func argfmt(driver string, nv driver.NamedValue) string {
	switch driver {
	// TODO: support more sql engines
	case "postgres":
		return fmt.Sprintf("$%d", nv.Ordinal)
	default:
		return "?"
	}
}
