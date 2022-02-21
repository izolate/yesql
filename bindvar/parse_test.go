package bindvar

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	tcs := []struct {
		driver string        // DB driver
		qt     string        // query template
		data   interface{}   // data object for args
		q      string        // returned query
		args   []interface{} // positional arg values
	}{
		{
			driver: "postgres",
			qt:     "SELECT * FROM a WHERE name = @Name AND age > @Age LIMIT 5",
			data: struct {
				Name string
				Age  int
			}{Name: "Foo", Age: 10},
			q:    "SELECT * FROM a WHERE name = $1 AND age > $2 LIMIT 5",
			args: []interface{}{"Foo", 10},
		},
		{
			driver: "postgres",
			qt:     "SELECT * FROM users WHERE name = @Name AND username LIKE '@foo%' LIMIT @Limit OFFSET 5",
			data: struct {
				Name  string
				Limit int
			}{Name: "Foo", Limit: 100},
			q:    "SELECT * FROM users WHERE name = $1 AND username LIKE '@foo%' LIMIT $2 OFFSET 5",
			args: []interface{}{"Foo", 100},
		},
		{
			driver: "postgres",
			qt:     `SELECT * FROM docs WHERE type = @Type AND dump = '{"text":"''@ignore @at @signs@@@''"}' LIMIT @Limit`,
			data: struct {
				Type  string
				Limit int
			}{Type: "Foo", Limit: 100},
			q:    `SELECT * FROM docs WHERE type = $1 AND dump = '{"text":"''@ignore @at @signs@@@''"}' LIMIT $2`,
			args: []interface{}{"Foo", 100},
		},
		{
			driver: "postgres",
			qt:     "SELECT * FROM strings WHERE locale = @Locale AND text ILIKE @J文 LIMIT @Limit",
			data: struct {
				Locale string
				J文     string
				Limit  int
			}{Locale: "JP", J文: "すみません", Limit: 3},
			q:    "SELECT * FROM strings WHERE locale = $1 AND text ILIKE $2 LIMIT $3",
			args: []interface{}{"JP", "すみません", 3},
		},
		{
			driver: "postgres",
			qt:     "SELECT created_at::timestamp(0) WHERE created_at > @date",
			data: map[string]interface{}{
				"date": time.Date(2020, 03, 10, 0, 0, 0, 0, time.UTC),
			},
			q:    "SELECT created_at::timestamp(0) WHERE created_at > $1",
			args: []interface{}{time.Date(2020, 03, 10, 0, 0, 0, 0, time.UTC)},
		},
		{
			driver: "postgres",
			qt:     "INSERT INTO authors (name) VALUES (@Name)",
			data: map[string]interface{}{
				"Name": "Max",
			},
			q:    "INSERT INTO authors (name) VALUES ($1)",
			args: []interface{}{"Max"},
		},
	}
	for _, tc := range tcs {
		bvar := New(tc.driver)
		q, args, err := bvar.Parse(tc.qt, tc.data)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if q != tc.q {
			t.Fatalf("Query not equal:\n%s\n-----\n%s\n", tc.q, q)
		}
		for i, a := range args {
			if a != tc.args[i] {
				t.Fatalf("Args not equal:\n%s\n-----\n%s\n", a, tc.args[i])
			}
		}
	}
}
