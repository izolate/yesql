package bindvar

import "fmt"

type Parser interface {
	// Parse parses all named parameters in a SQL statement, and returns
	// a statement with the params converted to bindvars appropriate for
	// the engine, e.g. :foo, :bar => $1, $2 (postgres).
	//
	// Additionally, the positional args are returned in order.
	Parse(query string, data interface{}) (q string, args []interface{}, err error)
}

func New() Parser {
	return &parser{}
}

type parser struct{}

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
		a  int // Pointer used to seek through the string
		ac int // Arg counter
	)
	for a < len(query) {
		// fmt.Printf("A:%d\n", a)

		ra := query[a] // the rune at position a

		// Find the previous and next runes to the current position.
		// We use these to accurately deduce a named parameter, and omit false positives.
		var ral, rar rune
		if a > 0 {
			ral = query[a-1]
		}
		if a < len(query)-1 {
			rar = query[a+1]
		}
		fmt.Printf("%s %s %s\n", string(ral), string(ra), string(rar))

		if string(ra) == ":" && string(ral) != ":" && string(rar) != ":" {
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
			a1 := a + 1 // Ignore prefixed colon (:) in arg name.
			args = append(args, []int{a1, b})

			// Add the positional placeholder to the output query
			// TODO: support more sql engines than postgres
			ac++ // increment arg counter
			p := []rune(fmt.Sprintf("$%d", ac))
			s = append(s, p...)

			// Skip to the end of arg
			a = b
			continue
		}

		s = append(s, query[a])
		a++
	}
	fmt.Println(string(s))
	return s, args
}
