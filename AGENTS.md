# Repository Guidelines

## Go Code Style

- Name every exported functional option constructor with the `Opt` prefix.
  This includes conditional options: use `OptQuietIf`, not `QuietIf`.
- Functional option constructors must return `func(*Config)` and follow the
  `OptX` naming pattern used by `OptDriver`, `OptTemplate`, and `OptQuiet`.
- Begin each exported option's doc comment with its exact function name.
