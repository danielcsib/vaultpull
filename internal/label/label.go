// Package label attaches arbitrary string metadata to secret maps.
package label

import "fmt"

// Set of labels keyed by name.
type Set map[string]string

// Labeler annotates secret maps with a fixed label set.
type Labeler struct {
	labels Set
	prefix string
}

// New creates a Labeler. prefix is prepended to every injected key.
func New(labels Set, prefix string) *Labeler {
	if prefix == "" {
		prefix = "VAULTPULL_"
	}
	return &Labeler{labels: labels, prefix: prefix}
}

// Apply injects labels into a copy of secrets as extra keys.
// Existing keys are never overwritten.
func (l *Labeler) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets)+len(l.labels))
	for k, v := range secrets {
		out[k] = v
	}
	for k, v := range l.labels {
		key := fmt.Sprintf("%s%s", l.prefix, k)
		if _, exists := out[key]; !exists {
			out[key] = v
		}
	}
	return out
}

// Strip removes all label-prefixed keys from a copy of secrets.
func (l *Labeler) Strip(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if len(k) < len(l.prefix) || k[:len(l.prefix)] != l.prefix {
			out[k] = v
		}
	}
	return out
}
