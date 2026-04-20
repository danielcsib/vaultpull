# plan

The `plan` package computes a **dry-run plan** of changes before any secrets are written to disk.

It compares the current state of an env file against the incoming secrets from Vault and produces a structured list of `Entry` values, each tagged with one of:

| Symbol | Action      | Meaning                        |
|--------|-------------|--------------------------------|
| `+`    | `add`       | Key is new in Vault            |
| `~`    | `update`    | Key exists but value changed   |
| `-`    | `delete`    | Key removed from Vault         |
| ` `    | `unchanged` | Key and value are identical    |

## Usage

```go
current := map[string]string{"DB_HOST": "localhost"}
incoming := map[string]string{"DB_HOST": "prod.db", "DB_PORT": "5432"}

p := plan.Build(".env", current, incoming)
if p.HasChanges() {
    p.Print(os.Stdout)
}
```

## Notes

- Entries are always sorted alphabetically by key for deterministic output.
- `HasChanges` returns `false` when every key is `unchanged`, making it safe to gate writes behind a plan check.
- Values are stored as-is; callers should apply masking before passing secrets if output will be shown to users.
