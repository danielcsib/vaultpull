# checkpoint

The `checkpoint` package tracks which Vault paths were last synced successfully, persisting metadata to a JSON file on disk.

## Usage

```go
store, err := checkpoint.NewStore(".vaultpull/checkpoint.json")
if err != nil {
    log.Fatal(err)
}

// After a successful sync:
store.Record("secret/myapp", len(secrets))

// Query last sync:
if entry, ok := store.Get("secret/myapp"); ok {
    fmt.Println("last synced:", entry.SyncedAt)
}
```

## Fields

| Field      | Description                          |
|------------|--------------------------------------|
| Path       | Vault secret path                    |
| SyncedAt   | UTC timestamp of last successful sync|
| KeyCount   | Number of keys written               |
