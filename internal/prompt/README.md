# prompt

Provides a simple interactive confirmation prompt for CLI commands.

## Usage

```go
c := prompt.New()
ok, err := c.Ask("Overwrite existing .env file?", false)
if err != nil {
    log.Fatal(err)
}
if !ok {
    fmt.Println("Aborted.")
    return
}
```

## Behaviour

| Input | defaultYes=false | defaultYes=true |
|-------|------------------|-----------------|
| `y` / `yes` | true | true |
| `n` / `no` | false | false |
| *(empty)* | false | true |
| EOF | false | true |

Any other input returns an error so the caller can re-prompt or abort.
