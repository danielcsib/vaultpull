# scope

The `scope` package provides root-prefix scoping for Vault secret paths.
It is useful when a sync configuration spans multiple paths but you want
to restrict or inspect only those paths that fall under a known prefix.

## Usage

```go
s, err := scope.New("secret/myapp")
if err != nil {
    log.Fatal(err)
}

// Check whether a single path is in scope.
if s.Contains("secret/myapp/db") {
    fmt.Println("in scope")
}

// Filter a slice of paths to only those within the root.
paths := []string{"secret/myapp/db", "secret/other", "secret/myapp/cache"}
scoped := s.Filter(paths)
// scoped == ["secret/myapp/db", "secret/myapp/cache"]

// Strip the root prefix to get a relative path.
rel := s.RelativePath("secret/myapp/db/credentials")
// rel == "db/credentials"
```

## API

| Function / Method | Description |
|---|---|
| `New(root string)` | Create a new Scoper; returns error if root is empty |
| `Contains(path string) bool` | Reports whether path is within the root |
| `Filter(paths []string) []string` | Returns only in-scope paths |
| `RelativePath(path string) string` | Strips root prefix from path |
| `Root() string` | Returns the normalised root prefix |
