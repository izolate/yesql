package yesql

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"reflect"
	"testing"
)

func TestLogStatement(t *testing.T) {
	var output bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&output, nil)))
	t.Cleanup(func() {
		slog.SetDefault(previous)
	})

	query := "SELECT *\nFROM books\nWHERE id = $1"
	args := []any{42}
	logStatement(context.Background(), false, query, args)

	var event struct {
		Level string `json:"level"`
		Msg   string `json:"msg"`
		Query string `json:"query"`
		Args  []any  `json:"args"`
	}
	if err := json.NewDecoder(&output).Decode(&event); err != nil {
		t.Fatalf("decode log event: %v", err)
	}
	if event.Level != slog.LevelInfo.String() {
		t.Errorf("level = %q; want %q", event.Level, slog.LevelInfo.String())
	}
	if event.Msg != "executing SQL statement" {
		t.Errorf("msg = %q; want %q", event.Msg, "executing SQL statement")
	}
	if event.Query != query {
		t.Errorf("query = %q; want %q", event.Query, query)
	}
	if want := []any{float64(42)}; !reflect.DeepEqual(event.Args, want) {
		t.Errorf("args = %#v; want %#v", event.Args, want)
	}

	output.Reset()
	logStatement(context.Background(), true, query, args)
	if output.Len() != 0 {
		t.Errorf("quiet log output = %q; want no output", output.String())
	}
}
