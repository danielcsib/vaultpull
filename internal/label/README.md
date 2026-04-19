# label

The `label` package injects metadata labels into secret maps as extra environment keys.

## Usage

```go
l := label.New(label.Set{
    "ENV":  "production",
    "TEAM": "platform",
}, "")

annotated := l.Apply(secrets)
// VAULTPULL_ENV=production
// VAULTPULL_TEAM=platform
```

## Behaviour

- Labels are injected with a configurable prefix (default: `VAULTPULL_`).
- Existing keys are **never** overwritten.
- `Apply` returns a new map; the input is not mutated.
- `Strip` removes all label-prefixed keys from a map, useful before writing the final `.env` file.
