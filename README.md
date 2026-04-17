# vaultpull

> CLI tool to sync secrets from HashiCorp Vault into local `.env` files safely

---

## Installation

```bash
go install github.com/yourusername/vaultpull@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpull/releases).

---

## Usage

Set your Vault address and token, then run `vaultpull` pointing at a secret path:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.xxxxxxxxxxxxxxxx"

vaultpull --path secret/data/myapp --output .env
```

This will fetch all key-value pairs stored at the given Vault path and write them to `.env`:

```
DB_HOST=db.example.com
DB_PASSWORD=supersecret
API_KEY=abc123
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--path` | Vault secret path to read from | *(required)* |
| `--output` | Output `.env` file path | `.env` |
| `--append` | Append to existing file instead of overwriting | `false` |
| `--dry-run` | Print secrets to stdout without writing to disk | `false` |
| `--prefix` | Add a prefix to all exported variable names | *(none)* |

### Example with dry run

```bash
vaultpull --path secret/data/myapp --dry-run
```

### Example with prefix

```bash
vaultpull --path secret/data/myapp --prefix MYAPP_
```

This would produce output like:

```
MYAPP_DB_HOST=db.example.com
MYAPP_DB_PASSWORD=supersecret
MYAPP_API_KEY=abc123
```

---

## Requirements

- Go 1.21+
- A running HashiCorp Vault instance
- A valid `VAULT_TOKEN` or other supported auth method

---

## License

[MIT](LICENSE) © 2024 yourusername
