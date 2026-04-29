# coalesce

The `coalesce` package merges multiple secret maps with a defined priority
order, returning the first non-empty value for each key across the provided
sources.

## Usage

```go
m := coalesce.New()

vaultSecrets := map[string]string{"DB_PASS": "vault-secret"}
defaults     := map[string]string{"DB_PASS": "default", "LOG_LEVEL": "info"}

out, err := m.Merge(vaultSecrets, defaults)
// out => {"DB_PASS": "vault-secret", "LOG_LEVEL": "info"}
```

## OmitEmpty

By default an empty string value from a higher-priority source wins over a
non-empty value in a lower-priority source. Use `OmitEmpty()` to skip empty
values and fall through to the next source:

```go
m := coalesce.New(coalesce.OmitEmpty())
```

## First

`First` is a convenience function that returns the first value found for a
given key across any number of source maps:

```go
val, ok := coalesce.First("API_KEY", overrides, base)
```

## Integration

`coalesce` is used by the sync pipeline to blend Vault secrets with local
defaults before writing the final `.env` file.
