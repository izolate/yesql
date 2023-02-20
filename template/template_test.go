package template

import "testing"

func TestExecTemplate(t *testing.T) {
	tcs := []struct {
		input  string
		data   any
		output string
	}{
		{
			input: `SELECT * FROM a WHERE name = :Name{{if .Status}} AND status = :Status{{end}}`,
			data: struct {
				Name   string
				Status string
			}{Name: "Foo"},
			output: `SELECT * FROM a WHERE name = :Name`,
		},
		{
			input: `SELECT * FROM a WHERE name = :Name{{if .Status}} AND status = :Status{{end}}`,
			data: struct {
				Name   string
				Status string
			}{Name: "bar", Status: "bar"},
			output: `SELECT * FROM a WHERE name = :Name AND status = :Status`,
		},
		{
			input: "SELECT p.id, p.name {{if .Status}}, s.status {{end}}FROM people p {{if .Status}}JOIN people_statuses s ON s.person_id = p.id {{end}}WHERE p.created_at > '2005-08-14' {{if .Status}}AND s.status = :Status {{end}}{{if .Limit}}LIMIT :limit {{end}}{{if .Offset}}OFFSET :offset {{end}}",
			data: struct {
				Status string
				Limit  int
				Offset int
			}{},
			output: "SELECT p.id, p.name FROM people p WHERE p.created_at > '2005-08-14' ",
		},
		{
			input: "SELECT p.id, p.name {{if .Status}}, s.status {{end}}FROM people p {{if .Status}}JOIN people_statuses s ON s.person_id = p.id {{end}}WHERE p.created_at > '2005-08-14' {{if .Status}}AND s.status = :Status {{end}}{{if .Limit}}LIMIT :limit {{end}}{{if .Offset}}OFFSET :offset {{end}}",
			data: struct {
				Status string
				Limit  int
				Offset int
			}{Status: "pending", Limit: 100},
			output: "SELECT p.id, p.name , s.status FROM people p JOIN people_statuses s ON s.person_id = p.id WHERE p.created_at > '2005-08-14' AND s.status = :Status LIMIT :limit ",
		},
		{
			input: "SELECT p.id, p.name {{if .Status}}, s.status {{end}}FROM people p {{if .Status}}JOIN people_statuses s ON s.person_id = p.id {{end}}WHERE p.created_at > '2005-08-14' {{if .Status}}AND s.status = :Status {{end}}{{if .Limit}}LIMIT :limit {{end}}{{if .Offset}}OFFSET :offset {{end}}",
			data: struct {
				Status string
				Limit  int
				Offset int
			}{Limit: 100, Offset: 5},
			output: "SELECT p.id, p.name FROM people p WHERE p.created_at > '2005-08-14' LIMIT :limit OFFSET :offset ",
		},
	}
	tpl := New()
	for _, tc := range tcs {
		result, err := tpl.Execute(tc.input, tc.data)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if result != tc.output {
			t.Fatalf("Not equal:\n%s\n-----\n%s\n", tc.output, result)
		}
	}
}
