# notify

Package `notify` provides a lightweight event notifier for vaultpull sync operations.

## Usage

```go
import "github.com/yourorg/vaultpull/internal/notify"

n := notify.New(os.Stdout)
n.Info("secret/app", "synced 5 keys")
n.Warn("secret/db", "no changes detected")
n.Error("secret/api", "vault unreachable")
```

## Event Format

Each event is written as a single line:

```
2024-01-15T10:00:00Z [INFO] path=secret/app synced 5 keys
```

## Levels

| Level   | Use case                        |
|---------|---------------------------------|
| INFO    | Successful sync, key counts     |
| WARN    | Skipped paths, stale cache      |
| ERROR   | Vault errors, write failures    |

## Notes

- If no writer is provided, `os.Stdout` is used.
- If no timestamp is set on the event, `time.Now()` is used automatically.
- Empty paths are rendered as `-` for log consistency.
