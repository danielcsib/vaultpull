# retry

Provides a simple exponential-backoff retry helper used when communicating with Vault.

## Usage

```go
p := retry.DefaultPolicy() // 3 attempts, 500ms base delay, 2x multiplier

err := retry.Do(ctx, p, func() error {
    return vaultClient.ReadSecrets(path)
})
if errors.Is(err, retry.ErrExhausted) {
    // all attempts failed
}
```

## Policy fields

| Field | Default | Description |
|---|---|---|
| MaxAttempts | 3 | Total number of tries |
| Delay | 500ms | Initial wait between attempts |
| Multiplier | 2.0 | Back-off growth factor |

Context cancellation is respected between attempts.
