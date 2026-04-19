# dedupe

The `dedupe` package merges multiple secret maps from different Vault paths, resolving key collisions according to a configurable policy.

## Policies

| Policy | Behaviour |
|---|---|
| `PolicyKeepFirst` | Retains the value from the first source seen |
| `PolicyKeepLast` | Overwrites with the value from the latest source |
| `PolicyError` | Returns an error immediately on any collision |

## Usage

```go
m := dedupe.New(dedupe.PolicyKeepLast)

base := map[string]string{"DB_HOST": "localhost"}
override := map[string]string{"DB_HOST": "prod.db", "API_KEY": "secret"}

if err := m.Merge(base, override); err != nil {
    log.Fatal(err)
}
// base now contains merged result
fmt.Println(m.Conflicts) // [{Key:"DB_HOST", ...}]
```

## Notes

- `Merge` mutates the `dst` map in place.
- Call `Reset()` between sync runs to clear conflict history.
