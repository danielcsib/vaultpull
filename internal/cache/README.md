# cache

The `cache` package provides a local on-disk cache of Vault secret snapshots.

## Purpose

Avoiding unnecessary writes to `.env` files when secrets have not changed
reduces noise in file-watchers and prevents accidental reloads of dependent
processes.

## How it works

- Each Vault path is hashed (SHA-256) to produce a stable filename.
- The cache file stores the secrets, their SHA-256 content hash, and a
  timestamp.
- `Changed()` compares the current secrets hash against the cached hash.

## Usage

```go
store, err := cache.NewStore(".vaultpull/cache")
if err != nil { ... }

changed, err := store.Changed("secret/myapp", fetchedSecrets)
if err != nil { ... }

if changed {
    // write env file
    _ = store.Set("secret/myapp", fetchedSecrets)
}
```

## Security

Cache files are written with mode `0600` and the cache directory is created
with mode `0700`. The directory should be added to `.gitignore`.
