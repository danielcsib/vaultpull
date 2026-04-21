# output

The `output` package serialises a `map[string]string` of secrets into a
chosen text format and writes it to any `io.Writer`.

## Supported Formats

| Format   | Description                                      |
|----------|--------------------------------------------------|
| `env`    | Plain `KEY=VALUE` lines, suitable for `.env` files |
| `export` | Shell-compatible `export KEY=VALUE` lines        |
| `json`   | Indented JSON object with sorted keys            |

## Usage

```go
import "github.com/your-org/vaultpull/internal/output"

secrets := map[string]string{
    "DB_HOST": "localhost",
    "DB_PORT": "5432",
}

// Write to stdout as plain env
w := output.New(output.FormatEnv, nil)
w.Write(secrets)

// Write to a buffer as JSON
var buf bytes.Buffer
w = output.New(output.FormatJSON, &buf)
w.Write(secrets)
```

## Notes

- Keys are always sorted alphabetically for deterministic output.
- Values containing spaces, tabs, newlines, or `#` are automatically quoted
  when using the `env` or `export` formats.
- The `json` format uses Go's `%q` verb for both keys and values.
