# expire

The `expire` package tracks TTL-based expiry for synced Vault secret paths.

## Usage

```go
store, err := expire.NewStore(".vaultpull-expire.json")
if err != nil {
    log.Fatal(err)
}

// Record that a path should be refreshed after 1 hour.
_ = store.Set("secret/app", time.Hour)

// Check before syncing.
if store.Expired("secret/app") {
    // re-sync
}

// Remove when no longer needed.
_ = store.Delete("secret/app")
```

## Behaviour

- Entries are persisted to a JSON file with mode `0600`.
- `Expired` returns `false` for unknown paths (no entry = never expires).
- A negative TTL can be used in tests to simulate an already-expired entry.
