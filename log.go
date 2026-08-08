package yesql

import (
	"context"
	"log/slog"
)

func logStatement(ctx context.Context, quiet bool, query string, args []any) {
	if quiet {
		return
	}
	slog.InfoContext(
		ctx,
		"executing SQL statement",
		slog.String("query", query),
		slog.Any("args", args),
	)
}
