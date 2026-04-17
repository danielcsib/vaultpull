# redact

The `redact` package provides value masking for secret data before it is
logged, printed, or included in reports.

## Usage

```go
r := redact.New(4) // reveal last 4 characters

masked := r.Mask("supersecret")
// => "*******cret"

maskedMap := r.MaskMap(map[string]string{
    "DB_PASS": "hunter2",
    "API_KEY": "abc123",
})
```

## Behaviour

- Empty strings are returned unchanged.
- If the value is shorter than or equal to `visibleChars`, the entire value is masked.
- `visibleChars = 0` masks the full value regardless of length.
- `MaskMap` returns a new map and does not modify the original.
