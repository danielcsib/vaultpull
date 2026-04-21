# quota

The `quota` package enforces per-path read limits to prevent runaway or accidental
bulk reads against HashiCorp Vault.

## Usage

```go
store := quota.New(5) // allow at most 5 reads per path

if err := store.Check("secret/data/myapp"); err != nil {
    // handle quota.ErrQuotaExceeded
}

fmt.Println(store.Remaining("secret/data/myapp")) // 4
```

## API

| Function | Description |
|---|---|
| `New(maxReads int) *Store` | Create a new quota store. Defaults to 10 if maxReads ≤ 0. |
| `Check(path string) error` | Increment counter; return `ErrQuotaExceeded` if over limit. |
| `Remaining(path string) int` | Returns reads left for the path (never negative). |
| `Reset(path string)` | Clear the counter for a single path. |
| `ResetAll()` | Clear all counters. |
| `Snapshot() map[string]int` | Return a copy of all current read counts. |
| `PrintSummary(s *Store, w io.Writer)` | Print a formatted usage table. |

## Notes

- All methods are safe for concurrent use.
- `Snapshot` returns a copy; mutations do not affect the store.
- `ErrQuotaExceeded` can be detected with `errors.Is`.
