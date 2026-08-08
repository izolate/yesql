package yesql

import (
	"context"
	"log/slog"
	"strings"
	"unicode"
)

func (c *Config) logSQL(ctx context.Context, query string) {
	if c.quiet {
		return
	}
	slog.InfoContext(
		ctx,
		"executing SQL statement",
		slog.String("query", inline(query)),
	)
}

// inline collapses whitespace outside SQL string literals.
func inline(query string) string {
	var b strings.Builder
	var quoted, space bool

	for _, r := range query {
		if r == '\'' {
			if space && b.Len() > 0 {
				b.WriteByte(' ')
			}
			space = false
			quoted = !quoted
			b.WriteRune(r)
			continue
		}
		if !quoted && unicode.IsSpace(r) {
			space = b.Len() > 0
			continue
		}
		if space {
			b.WriteByte(' ')
			space = false
		}
		b.WriteRune(r)
	}

	return b.String()
}
