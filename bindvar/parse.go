package bindvar

import (
	"database/sql/driver"
	"fmt"
	"reflect"
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

func New(driver string) Parser {
	return &parser{driver}
}

type parser struct {
	driver string
}

func (parser) Parse(query string, data interface{}) (string, []interface{}, error) {
	// Convert to rune to handle unicode strings
	qt := []rune(query)

	// Parse named args
	q, nvs := parse(qt)
	args := []interface{}{}
	for _, nv := range nvs {
		// Get the named arg values from data
		v := value(data, nv.Name)
		args = append(args, v)
		fmt.Printf("%d) %v %v\n", nv.Ordinal, nv.Name, v)
	}

	return string(q), args, nil
}

func parse(query []rune) (s []rune, args []driver.NamedValue) {
	var (
		a      int  // Pointer used to seek through the string
		op     int  // The ordinal position of the captured arg
		ignore bool // Used to ignore false positives
	)
	for a < len(query) {
		ra := query[a] // the rune at position a

		fmt.Printf("%s\n", string(ra))

		// Ignore characters inside sql string literals.
		if string(ra) == "'" {
			ignore = !ignore
		}

		if !ignore && string(ra) == naPrefix {
			// We've found an argument! Create second pointer to find end of argument.
			b := a

			// Find the first non-allowed rune to infer the end of the arg.
			for b < len(query) {
				fmt.Printf("B:%d\n", b)
				rb := query[b]
				if string(rb) == " " {
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
			// TODO: support more sql engines than postgres
			arg := []rune(formatArg("postgres", nv))
			s = append(s, arg...)

			a = b // Skip to the end of the arg
			continue
		}

		s = append(s, query[a])
		a++
	}
	return s, args
}

func value(data interface{}, name string) interface{} {
	if m, ok := data.(map[string]interface{}); ok {
		if v, ok := m[name]; ok {
			return v
		}
		return nil
	}
	if v := reflect.Indirect(reflect.ValueOf(data)); v.Kind() == reflect.Struct {
		if f := v.FieldByName(name); f.IsValid() {
			return f.Interface()
		}
	}
	return nil
}

func formatArg(driver string, nv driver.NamedValue) string {
	switch driver {
	case "postgres":
		return fmt.Sprintf("$%d", nv.Ordinal)
	default:
		return "?"
	}
}
