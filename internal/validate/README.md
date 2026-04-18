# validate

The `validate` package checks a resolved secret map against a set of required keys before writing to disk.

## Usage

```go
import "github.com/your-org/vaultpull/internal/validate"

result := validate.Check(secrets, []string{"DB_URL", "API_KEY", "JWT_SECRET"})
if result.HasErrors() {
    fmt.Print(result.Summary())
    os.Exit(1)
}
```

## Checks performed

| Check | Severity | Description |
|-------|----------|-------------|
| Missing key | Error | A required key was not found in the secret map |
| Empty value | Warning | Key exists but its value is an empty string |
| Whitespace | Warning | Value has leading or trailing whitespace |

## Result

`Result.HasErrors()` returns `true` only when required keys are absent — empty values and whitespace are non-fatal warnings surfaced via `Summary()`.
