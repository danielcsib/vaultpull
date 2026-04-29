# internal/drift

Package `drift` detects when local `.env` files have diverged from the secrets stored in Vault.

## Overview

After a `vaultpull sync`, local files should mirror Vault exactly. Over time they can drift due to manual edits, out-of-band deployments, or key rotations. This package compares the two maps and produces a structured `Report`.

## Usage

```go
vaultSecrets := map[string]string{"DB_PASS": "new-secret", "API_KEY": "abc"}
localEnv    := map[string]string{"DB_PASS": "old-secret"}

report, err := drift.Detect("secret/myapp", vaultSecrets, localEnv)
if err != nil {
    log.Fatal(err)
}

if report.HasDrift() {
    for _, e := range report.Drifted() {
        fmt.Printf("%-30s %s\n", e.Key, e.Status)
    }
}
```

## Statuses

| Status    | Meaning                                          |
|-----------|--------------------------------------------------|
| `match`   | Local value equals Vault value                   |
| `drifted` | Local value differs from Vault value             |
| `missing` | Key exists in Vault but is absent locally        |
| `orphan`  | Key exists locally but is not present in Vault   |

## Notes

- `Detect` treats Vault as the source of truth.
- Entries in the returned `Report` are sorted alphabetically by key.
- An empty `path` argument returns an error.
