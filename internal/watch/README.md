# internal/watch

Polls one or more Vault secret paths at a fixed interval and invokes a
callback whenever the returned secrets differ from the previously seen values.

## Usage

```go
client, _ := vault.NewClient(vault.Config{...})

w := watch.New(client, []string{"secret/app", "secret/db"}, 30*time.Second,
    func(path string, secrets map[string]string) {
        fmt.Printf("secrets changed at %s\n", path)
    },
)

ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
defer cancel()
w.Run(ctx) // blocks until ctx is cancelled
```

## Behaviour

- On the **first poll** the callback is always fired (no previous state).
- Subsequent polls only fire the callback when at least one key/value pair
  differs from the last observed snapshot.
- Read errors are logged and skipped; the previous snapshot is retained.
- A zero or negative interval falls back to **30 seconds**.
