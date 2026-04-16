# audit

The `audit` package provides append-only JSON Lines audit logging for `vaultpull` sync operations.

## Usage

```go
logger := audit.NewLogger("/var/log/vaultpull/audit.jsonl")

logger.Record(audit.Entry{
    SecretPath: "secret/data/myapp",
    EnvFile:    ".env.production",
    Keys:       []string{"DB_PASSWORD", "API_KEY"},
    Success:    true,
})
```

## Format

Each sync operation appends one JSON object per line:

```json
{"timestamp":"2024-01-15T10:30:00Z","secret_path":"secret/data/myapp","env_file":".env.production","keys":["DB_PASSWORD","API_KEY"],"success":true}
```

Failed syncs include an `error` field:

```json
{"timestamp":"2024-01-15T10:31:00Z","secret_path":"secret/data/other","env_file":".env","keys":null,"success":false,"error":"permission denied"}
```

## Configuration

Set `audit_log` in your `vaultpull.yaml` to enable logging. An empty path disables audit logging entirely.
