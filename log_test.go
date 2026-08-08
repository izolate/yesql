package yesql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"testing"
)

func TestLogSQL(t *testing.T) {
	var output bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&output, nil)))
	t.Cleanup(func() {
		slog.SetDefault(previous)
	})

	query := "SELECT *\nFROM books\nWHERE id = $1"
	new(Config).logSQL(context.Background(), query)

	var event map[string]any
	if err := json.NewDecoder(&output).Decode(&event); err != nil {
		t.Fatalf("decode log event: %v", err)
	}
	if got := event["level"]; got != slog.LevelInfo.String() {
		t.Errorf("level = %q; want %q", got, slog.LevelInfo.String())
	}
	if got := event["msg"]; got != "executing SQL statement" {
		t.Errorf("msg = %q; want %q", got, "executing SQL statement")
	}
	if got, want := event["query"], inline(query); got != want {
		t.Errorf("query = %q; want %q", got, want)
	}
	if _, ok := event["args"]; ok {
		t.Error("log event includes bound arguments")
	}

	output.Reset()
	(&Config{quiet: true}).logSQL(context.Background(), query)
	if output.Len() != 0 {
		t.Errorf("quiet log output = %q; want no output", output.String())
	}
}

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
			in: `

				SELECT

				id

				FROM

				foo
			`,
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
			if got := inline(tc.in); got != tc.out {
				t.Errorf("inline() = %q; want %q", got, tc.out)
			}
		})
	}
}
