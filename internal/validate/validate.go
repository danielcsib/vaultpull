package validate

import (
	"fmt"
	"strings"
)

// Result holds validation findings for a secret map.
type Result struct {
	Missing  []string
	Empty    []string
	Warnings []string
}

// HasErrors returns true if any required keys are missing.
func (r Result) HasErrors() bool {
	return len(r.Missing) > 0
}

// Summary returns a human-readable summary string.
func (r Result) Summary() string {
	var sb strings.Builder
	for _, k := range r.Missing {
		sb.WriteString(fmt.Sprintf("  [ERROR] missing required key: %s\n", k))
	}
	for _, k := range r.Empty {
		sb.WriteString(fmt.Sprintf("  [WARN]  key has empty value: %s\n", k))
	}
	for _, w := range r.Warnings {
		sb.WriteString(fmt.Sprintf("  [WARN]  %s\n", w))
	}
	return sb.String()
}

// Check validates secrets against a list of required keys.
// It reports missing keys, empty values, and keys with suspicious whitespace.
func Check(secrets map[string]string, required []string) Result {
	var r Result

	for _, key := range required {
		val, ok := secrets[key]
		if !ok {
			r.Missing = append(r.Missing, key)
			continue
		}
		if val == "" {
			r.Empty = append(r.Empty, key)
		}
	}

	for k, v := range secrets {
		if strings.TrimSpace(v) != v && v != "" {
			r.Warnings = append(r.Warnings, fmt.Sprintf("key %q has leading/trailing whitespace", k))
		}
	}

	return r
}
