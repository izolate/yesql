package bindvar

import (
	"database/sql/driver"
	"fmt"
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

	//
	q, argn := parse(qt)
	for _, a := range argn {
		s := qt[a[0]:a[1]]
		fmt.Println(string(s))
	}

	return string(q), nil, nil
}

func parse(query []rune) (s []rune, args [][]int) {
	var (
		a      int  // Pointer used to seek through the string
		ac     int  // Total number of captured args
		ignore bool // Used to ignore false positives
	)
	for a < len(query) {
		// fmt.Printf("A:%d\n", a)

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

			// Collect the placeholder names to later extract values from data.
			a1 := a + 1 // Ignore prefix (@) in arg name.
			args = append(args, []int{a1, b})

			// Convert the named arg to the correct arg for the driver.
			// TODO: support more sql engines than postgres
			ac++ // increment arg counter
			nv := driver.NamedValue{
				Ordinal: ac,
				Name:    string(query[a1:b]),
			}
			arg := []rune(formatArg("postgres", nv))
			s = append(s, arg...)

			// Skip to the end of arg.
			a = b
			continue
		}

		s = append(s, query[a])
		a++
	}
	fmt.Println(string(s))
	return s, args
}

func formatArg(driver string, nv driver.NamedValue) string {
	switch driver {
	case "postgres":
		return fmt.Sprintf("$%d", nv.Ordinal)
	default:
		return "?"
	}
}
