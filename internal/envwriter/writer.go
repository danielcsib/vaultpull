package envwriter

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// WriteEnvFile writes a map of key-value secret pairs to a .env file at the given path.
// Existing file contents are overwritten. Keys are written in sorted order for determinism.
func WriteEnvFile(path string, secrets map[string]string) error {
	if path == "" {
		return fmt.Errorf("envwriter: output path must not be empty")
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		v := secrets[k]
		sb.WriteString(formatLine(k, v))
		sb.WriteByte('\n')
	}

	// Write with restrictive permissions — secrets should not be world-readable.
	if err := os.WriteFile(path, []byte(sb.String()), 0600); err != nil {
		return fmt.Errorf("envwriter: failed to write file %q: %w", path, err)
	}
	return nil
}

// formatLine formats a single key=value pair, quoting the value if it contains
// spaces, newlines, or special shell characters.
func formatLine(key, value string) string {
	if needsQuoting(value) {
		escaped := strings.ReplaceAll(value, `"`, `\"`)
		return fmt.Sprintf(`%s="%s"`, key, escaped)
	}
	return fmt.Sprintf("%s=%s", key, value)
}

func needsQuoting(v string) bool {
	for _, ch := range v {
		switch ch {
		case ' ', '\t', '\n', '\r', '"', '\'', '\\', '#', '$', '`', '!':
			return true
		}
	}
	return false
}

// ValidateKeys checks that all keys in the map are non-empty and contain only
// alphanumeric characters or underscores, as required by most shell environments.
func ValidateKeys(secrets map[string]string) error {
	for k := range secrets {
		if k == "" {
			return fmt.Errorf("envwriter: secret key must not be empty")
		}
		for _, ch := range k {
			if !isValidKeyChar(ch) {
				return fmt.Errorf("envwriter: invalid character %q in key %q", ch, k)
			}
		}
	}
	return nil
}

func isValidKeyChar(ch rune) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '_'
}
