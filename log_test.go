package yesql

import (
	"fmt"
	"testing"
)

func TestInline(t *testing.T) {
	testCases := []struct {
		in  string
		out string
	}{
		{
			in:  `SELECT * FROM foo WHERE id = $1`,
			out: `SELECT * FROM foo WHERE id = $1`,
		},
		{
			in: `SELECT *
				FROM foo
				WHERE id = $1`,
			out: `SELECT * FROM foo WHERE id = $1`,
		},
		{
			in: `SELECT

				id

				FROM

				foo`,
			out: `SELECT id FROM foo`,
		},
		{
			in: `INSERT INTO documents (
				id,
				value
			)
			VALUES (
				1,
				'{"name": " ignore    spaces    here   "}'::jsonb
			)`,
			out: `INSERT INTO documents ( id, value ) VALUES ( 1, '{"name": " ignore    spaces    here   "}'::jsonb )`,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			its := assert{t}
			its.StringEq(tc.out, inline(tc.in))
		})
	}
}
