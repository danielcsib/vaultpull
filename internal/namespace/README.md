# namespace

The `namespace` package scopes secret keys under a configurable prefix,
enabling clean separation of secrets across environments or services.

## Usage

```go
s, err := namespace.New("prod", "/")
if err != nil {
    log.Fatal(err)
}

// Qualify individual keys
key := s.Qualify("DB_PASSWORD") // "prod/DB_PASSWORD"

// Apply namespace to an entire map
scoped := s.Apply(map[string]string{
    "DB_HOST": "db.internal",
    "API_KEY": "secret",
})
// {"prod/DB_HOST": "db.internal", "prod/API_KEY": "secret"}

// Strip namespace from keys read back from Vault
plain := s.Unwrap(scoped)
// {"DB_HOST": "db.internal", "API_KEY": "secret"}
```

## Notes

- The separator defaults to `"/"` when an empty string is provided.
- `Apply` and `Unwrap` never mutate the input map.
- Keys that do not carry the expected prefix are passed through unchanged by `Strip` / `Unwrap`.
