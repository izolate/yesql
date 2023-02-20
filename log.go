package yesql

import (
	"log"
)

// a stop-gap measure until structured logging is added to the std lib.
// see: https://go.googlesource.com/proposal/+/master/design/56345-structured-logging.md
func logStatement(quiet bool, stmt string, args []any) {
	if quiet {
		return
	}
	log.SetFlags(log.LstdFlags)
	log.Println(inline(stmt))
}

func inline(s string) (t string) {
	var ignore bool
	for i, c := range s {
		// Ignore characters inside sql string literals.
		if string(c) == "'" {
			ignore = !ignore
		}
		first := i == 0
		if !ignore {
			switch {
			case !first && c == ' ' && t[len(t)-1] == ' ':
				// remove duplicate space
				continue

			case !first && c == '\n' && t[len(t)-1] != ' ':
				// convert the first newline to a space
				t += string(' ')
				continue

			case c == '\n', c == '\t':
				continue
			}
		}
		t += string(c)
	}
	return t
}
