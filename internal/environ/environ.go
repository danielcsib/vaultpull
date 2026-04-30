// Package environ provides utilities for overlaying secret maps onto
// the current process environment, producing a merged key-value set
// that can be passed to child processes or written to .env files.
package environ

import (
	"os"
	"strings"
)

// Options controls how the overlay behaves.
type Options struct {
	// Overwrite allows vault secrets to overwrite existing OS env vars.
	// When false, OS env vars take precedence.
	Overwrite bool

	// Prefix restricts which OS env vars are included in the base layer.
	// An empty prefix includes all OS env vars.
	Prefix string
}

// DefaultOptions returns a sensible default: secrets overwrite OS env,
// and all OS env vars are included.
func DefaultOptions() Options {
	return Options{Overwrite: true}
}

// Overlay merges secrets on top of (or beneath) the current OS environment
// and returns the combined map. The original secrets map is never mutated.
func Overlay(secrets map[string]string, opts Options) map[string]string {
	result := make(map[string]string)

	// Layer 1: OS environment (optionally filtered by prefix).
	for _, entry := range os.Environ() {
		k, v, ok := splitEntry(entry)
		if !ok {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(k, opts.Prefix) {
			continue
		}
		result[k] = v
	}

	// Layer 2: secrets.
	for k, v := range secrets {
		if _, exists := result[k]; exists && !opts.Overwrite {
			continue
		}
		result[k] = v
	}

	return result
}

// ToSlice converts a map into a slice of "KEY=VALUE" strings suitable
// for use with exec.Cmd.Env.
func ToSlice(env map[string]string) []string {
	out := make([]string, 0, len(env))
	for k, v := range env {
		out = append(out, k+"="+v)
	}
	return out
}

func splitEntry(entry string) (key, value string, ok bool) {
	idx := strings.IndexByte(entry, '=')
	if idx < 1 {
		return "", "", false
	}
	return entry[:idx], entry[idx+1:], true
}
