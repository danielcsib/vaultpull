# internal/tag

The `tag` package provides a lightweight tagging system for annotating secret keys with arbitrary metadata. Tags are key=value pairs that can be used to classify secrets by environment, team, sensitivity level, or any other dimension.

## Usage

```go
tr := tag.New()

// Attach tags to secret keys
tr.Set("DB_PASSWORD", "env", "prod")
tr.Set("DB_PASSWORD", "sensitivity", "high")
tr.Set("API_KEY", "env", "prod")

// Check for a specific tag
if tr.Has("DB_PASSWORD", "sensitivity", "high") {
    fmt.Println("DB_PASSWORD is highly sensitive")
}

// Filter a secrets map by tag constraints (AND semantics)
result := tr.Filter(secrets, []tag.Tag{
    {Key: "env", Value: "prod"},
    {Key: "sensitivity", Value: "high"},
})
```

## Behaviour

- Setting a tag with the same `tagKey` on the same `secretKey` **overwrites** the previous value.
- `Filter` applies constraints with **AND** semantics — all constraints must match.
- Empty `secretKey` or `tagKey` arguments are silently ignored.
- `String` returns a human-readable summary of all tags for a given key.
