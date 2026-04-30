# pin

The `pin` package lets operators lock a Vault secret path to a specific
version so that `vaultpull` will never silently overwrite it during a
routine sync.

## Concepts

| Term | Description |
|------|-------------|
| **Entry** | A record that associates a Vault path with a KV version number. |
| **Store** | A JSON-backed store that persists pin entries across runs. |

## Usage

```go
store, err := pin.NewStore(".vaultpull/pins.json")
if err != nil {
    log.Fatal(err)
}

// Pin secret/app/db to version 3.
if err := store.Pin("secret/app/db", 3, "alice"); err != nil {
    log.Fatal(err) // ErrAlreadyPinned if a different version is already pinned
}

// Check during sync.
if e, ok := store.Get("secret/app/db"); ok {
    fmt.Printf("path is pinned to v%d\n", e.Version)
}

// Remove pin.
_ = store.Unpin("secret/app/db")
```

## CLI integration

The syncer checks pinned paths before writing and skips any secret whose
current Vault version differs from the pinned version, logging a warning
instead of overwriting the local file.
