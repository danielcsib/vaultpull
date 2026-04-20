# internal/lineage

The `lineage` package records the provenance of every secret written to a
`.env` file: which Vault path and key each environment variable was sourced
from, and when it was fetched.

## Usage

```go
store, err := lineage.NewStore(".vaultpull/lineage.json")
if err != nil {
    log.Fatal(err)
}

// After fetching a secret from Vault:
store.Record("DB_PASSWORD", "secret/myapp/db", "password")

// Persist to disk:
if err := store.Flush(); err != nil {
    log.Fatal(err)
}

// Print a summary table:
lineage.PrintSummary(os.Stdout, store)
```

## API

| Symbol | Description |
|---|---|
| `NewStore(path)` | Load (or create) a lineage store at `path`; pass `""` for in-memory only |
| `(*Store).Record(envKey, vaultPath, vaultKey)` | Record provenance for an env key |
| `(*Store).Get(envKey)` | Retrieve the entry for a single key |
| `(*Store).All()` | Return all recorded entries |
| `(*Store).Flush()` | Persist the store to disk |
| `PrintSummary(w, store)` | Write a human-readable provenance table |

## Storage format

Entries are stored as a JSON object keyed by env variable name:

```json
{
  "DB_PASSWORD": {
    "env_key": "DB_PASSWORD",
    "vault_path": "secret/myapp/db",
    "vault_key": "password",
    "fetched_at": "2024-06-01T12:00:00Z"
  }
}
```
