# envlock

Package `envlock` provides advisory file-based locking for `.env` files,
preventing concurrent `vaultpull` processes from writing to the same file
simultaneously.

## How it works

A lock file (`<target>.lock`) is created alongside the target `.env` file
before any write operation begins and removed once the write completes.

Lock files carry an implicit TTL (default **30 s**). If a lock file is older
than the TTL it is considered **stale** and automatically replaced on the
next `Acquire` call, guarding against leftover locks from crashed processes.

## Usage

```go
l := envlock.New(".env", 0) // 0 → DefaultTTL (30 s)

if err := l.Acquire(); err != nil {
    if errors.Is(err, envlock.ErrLockHeld) {
        log.Fatal("another vaultpull process is already running")
    }
    log.Fatal(err)
}
defer l.Release()

// safe to write .env here
```

## API

| Symbol | Description |
|---|---|
| `New(path, ttl)` | Create a Locker for the given env file path |
| `Acquire()` | Obtain the lock; returns `ErrLockHeld` if already held |
| `Release()` | Release the lock (idempotent) |
| `Held()` | Report whether the lock is currently held |
| `ErrLockHeld` | Sentinel error returned by `Acquire` |
