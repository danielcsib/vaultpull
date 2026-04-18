// Package truncate provides utilities for truncating secret values
// to a maximum byte length before writing to .env files.
package truncate

import "fmt"

// DefaultMaxBytes is the default maximum byte length for a secret value.
const DefaultMaxBytes = 1024

// Truncator holds configuration for truncation behaviour.
type Truncator struct {
	maxBytes int
	suffix   string
}

// New returns a Truncator with the given max byte limit.
// If maxBytes is <= 0 the DefaultMaxBytes is used.
func New(maxBytes int, suffix string) *Truncator {
	if maxBytes <= 0 {
		maxBytes = DefaultMaxBytes
	}
	if suffix == "" {
		suffix = "..."
	}
	return &Truncator{maxBytes: maxBytes, suffix: suffix}
}

// Value truncates a single string value if it exceeds the byte limit.
func (t *Truncator) Value(v string) string {
	if len(v) <= t.maxBytes {
		return v
	}
	cutoff := t.maxBytes - len(t.suffix)
	if cutoff < 0 {
		cutoff = 0
	}
	return v[:cutoff] + t.suffix
}

// Apply truncates all values in the provided map, returning a new map.
// Keys that had values truncated are collected in the returned slice.
func (t *Truncator) Apply(secrets map[string]string) (map[string]string, []string) {
	out := make(map[string]string, len(secrets))
	var truncated []string
	for k, v := range secrets {
		result := t.Value(v)
		out[k] = result
		if result != v {
			truncated = append(truncated, fmt.Sprintf("%s (%d→%d bytes)", k, len(v), len(result)))
		}
	}
	return out, truncated
}
