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

// inline collapses whitespace outside single-quoted SQL string literals.
// It is intended only for readable logs, not as a complete SQL parser;
// dialect-specific quoting and SQL comments receive no special handling.
func inline(query string) string {
	var b strings.Builder
	var quoted, space bool

	for _, r := range query {
		// Preserve string literals exactly, including their whitespace.
		if r == '\'' {
			if space && b.Len() > 0 {
				b.WriteByte(' ')
			}
			space = false
			quoted = !quoted
			b.WriteRune(r)
			continue
		}
		// Defer whitespace until the next token. This collapses runs and
		// drops whitespace at the beginning and end of the query.
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
