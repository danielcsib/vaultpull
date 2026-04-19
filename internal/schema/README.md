# schema

The `schema` package validates that a map of secrets conforms to a declared set of rules before they are written to disk.

## Usage

```go
rules := []schema.FieldRule{
    {Key: "DB_HOST", Required: true},
    {Key: "DB_PASSWORD", Required: true, AllowEmpty: false},
    {Key: "LOG_LEVEL", Required: false, AllowEmpty: true},
}

s := schema.New(rules)
violations, err := s.Validate(secrets)
if err != nil {
    for _, v := range violations {
        fmt.Println(v)
    }
}
```

## FieldRule

| Field | Type | Description |
|-------|------|-------------|
| `Key` | string | Secret key name |
| `Required` | bool | Fail if key is absent |
| `AllowEmpty` | bool | Permit blank/whitespace values |

## Violation

Each `Violation` carries the offending `Key` and a human-readable `Message`. The slice is returned alongside a sentinel error when any rule is breached.
