package bindvar

import "testing"

func TestParse(t *testing.T) {
	tcs := []struct {
		qt   string        // query template
		data interface{}   // data object for args
		q    string        // returned query
		args []interface{} // positional arg values
	}{
		/*
			{
				qt: "SELECT * FROM a WHERE name = :Name AND age > :Age LIMIT 5",
				data: struct {
					Name string
					Age  int
				}{Name: "Foo", Age: 10},
				q:    "SELECT * FROM a WHERE name = $1 AND age > $2 LIMIT 5",
				args: []interface{}{"Foo", 10},
			},
		*/
		{
			qt: "SELECT * FROM urls WHERE title ILIKE :Title AND domain = 'http://example.com' LIMIT :Limit OFFSET 5",
			data: struct {
				Title string
				Limit int
			}{Title: "Foo", Limit: 100},
			q:    "SELECT * FROM urls WHERE title ILIKE $1 AND domain = 'http://example.com' LIMIT $2 OFFSET 5",
			args: []interface{}{"Foo", 100},
		},
		/*
			{
				qt: "SELECT * FROM strings WHERE locale = :Locale AND text ILIKE :很好 LIMIT :Limit",
				data: struct {
					Locale string
					很好     string
					Limit  int
				}{Locale: "CN", 很好: "很好", Limit: 3},
				q:    "SELECT * FROM strings WHERE locale = $1 AND text ILIKE $2 LIMIT $3",
				args: []interface{}{"Foo", "很好", 3},
			},
		*/
	}
	bvar := New()
	for _, tc := range tcs {
		q, _, err := bvar.Parse(tc.qt, tc.data)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if q != tc.q {
			t.Fatalf("Not equal:\n%s\n-----\n%s\n", tc.q, q)
		}
		// TODO: assert arg values
	}
}
