# sanitize

The `sanitize` package cleans and normalises secret values before they are
written to `.env` files, preventing hidden characters or unexpected whitespace
from causing hard-to-debug runtime issues.

## Features

- **TrimSpace** — strips leading and trailing whitespace from every value.
- **StripNonPrintable** — removes non-printable Unicode code points while
  preserving newlines (`\n`) and tabs (`\t`).
- **NormaliseNewlines** — converts Windows (`\r\n`) and old Mac (`\r`) line
  endings to Unix (`\n`).

## Usage

```go
import "github.com/yourusername/vaultpull/internal/sanitize"

// Use the recommended defaults.
s := sanitize.New(sanitize.DefaultOptions())

// Sanitise a single value.
clean := s.Value("  my secret\r\n")
// → "my secret"

// Sanitise an entire secrets map.
secrets := map[string]string{
    "DB_PASS": "  hunter2  ",
    "API_KEY": "abc\x00def",
}
clean := s.Apply(secrets)
// → {"DB_PASS": "hunter2", "API_KEY": "abcdef"}
```

## Options

| Field | Default | Description |
|---|---|---|
| `TrimSpace` | `true` | Remove surrounding whitespace |
| `StripNonPrintable` | `true` | Drop non-printable runes |
| `NormaliseNewlines` | `true` | Unify line endings to `\n` |
