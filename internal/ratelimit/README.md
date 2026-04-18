# ratelimit

Provides a simple token-bucket rate limiter for controlling the frequency of Vault API calls.

## Usage

```go
limiter, err := ratelimit.New(10) // 10 requests per second
if err != nil {
    log.Fatal(err)
}

for _, path := range paths {
    if err := limiter.Wait(ctx); err != nil {
        return err
    }
    secrets, err := client.ReadSecrets(ctx, path)
    // ...
}
```

## Behaviour

- Tokens are refilled based on elapsed wall-clock time.
- `Wait` blocks until a token is available or the context is cancelled.
- Safe for concurrent use.
