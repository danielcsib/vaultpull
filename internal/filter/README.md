# filter

The `filter` package provides prefix-based include/exclude rules for secret keys
read from Vault before they are written to `.env` files.

## Usage

```go
f := filter.New([]string{
    "APP_",          // include keys starting with APP_
    "!APP_INTERNAL_", // but exclude APP_INTERNAL_*
})

allowed := f.Apply(secrets)
```

## Rule Syntax

| Pattern       | Meaning                              |
|---------------|--------------------------------------|
| `APP_`        | Include keys with this prefix        |
| `!INTERNAL_`  | Exclude keys with this prefix        |

## Precedence

- Exclude rules (`!`) are evaluated **before** include rules.
- If **no include rules** are defined, all keys are allowed by default (unless excluded).
- If **at least one include rule** exists, only matching keys are included.
