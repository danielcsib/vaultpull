# merge

The `merge` package combines an existing `.env` file with a fresh set of secrets pulled from Vault.

## Behaviour

- **Preserved** — keys that exist only in the local file are kept as-is (e.g. developer overrides).
- **Added** — keys present in the incoming secrets but not in the file.
- **Updated** — keys present in both, where the incoming value differs from the file value.

Incoming secrets always win for overlapping keys.

## Usage

```go
result, err := merge.FromFile(".env", pulledSecrets)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("added: %v\n", result.Added)
fmt.Printf("updated: %v\n", result.Updated)
fmt.Printf("preserved: %v\n", result.Preserved)

// Write result.Final back to disk via envwriter.
```

## Notes

- Lines beginning with `#` are treated as comments and ignored.
- Quoted values (`KEY="val"`) are unquoted during parsing.
- If the target file does not exist, all incoming keys are treated as Added.
