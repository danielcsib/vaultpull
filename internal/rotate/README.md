# rotate

The `rotate` package provides backup rotation for `.env` files before they are overwritten by `vaultpull sync`.

## Usage

```go
r := rotate.New(".vaultpull/backups", 5)
if err := r.Rotate(".env"); err != nil {
    log.Fatal(err)
}
```

## Behaviour

- If the target `.env` file does not yet exist, `Rotate` is a no-op.
- Each call copies the current file to `<backupDir>/<basename>.<timestamp>.bak`.
- Timestamps use the format `20060102T150405Z` (UTC), ensuring lexicographic sort equals chronological sort.
- After writing, the oldest backups are pruned so that at most `maxBackups` copies are retained.
- Default `maxBackups` is **5** when a value ≤ 0 is supplied.

## Security

- Backup files are written with mode `0600`.
- The backup directory is created with mode `0700`.
