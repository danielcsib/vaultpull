# rollback

The `rollback` package restores `.env` files from previously created backups.

## Usage

```go
rb := rollback.New("/path/to/backups")

// Restore the latest backup for a given target file.
src, err := rb.Restore(".env")
if err != nil {
    log.Fatal(err)
}
fmt.Println("Restored from", src)

// Or restore from a specific backup.
err = rb.RestoreFrom("/path/to/backups/.env.20240101T120000.bak", ".env")
```

## Backup naming

Backups are expected to follow the naming convention produced by the `rotate`
package: `<basename>.<timestamp>.bak`. The `Latest` method selects the
lexicographically last match, which corresponds to the most recent timestamp.

## Notes

- `Restore` overwrites the target file in place.
- No confirmation is requested; use `internal/prompt` if interactive confirmation is needed.
