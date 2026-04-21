// Package namespace provides utilities for scoping secret keys
// under a named namespace prefix, enabling multi-environment
// or multi-service secret isolation.
package namespace

import (
	"errors"
	"fmt"
	"strings"
)

// Scoper applies and strips namespace prefixes from secret maps.
type Scoper struct {
	prefix    string
	separator string
}

// New creates a Scoper with the given namespace and separator.
// The separator defaults to "/" if empty.
func New(ns, separator string) (*Scoper, error) {
	if strings.TrimSpace(ns) == "" {
		return nil, errors.New("namespace: namespace must not be empty")
	}
	sep := separator
	if sep == "" {
		sep = "/"
	}
	return &Scoper{prefix: ns, separator: sep}, nil
}

// Qualify returns the fully-qualified key: "<namespace><sep><key>".
func (s *Scoper) Qualify(key string) string {
	return fmt.Sprintf("%s%s%s", s.prefix, s.separator, key)
}

// Strip removes the namespace prefix from key.
// If the key does not start with the prefix, it is returned unchanged.
func (s *Scoper) Strip(key string) string {
	expected := s.prefix + s.separator
	if strings.HasPrefix(key, expected) {
		return strings.TrimPrefix(key, expected)
	}
	return key
}

// Apply returns a new map with every key qualified under the namespace.
// The original map is not mutated.
func (s *Scoper) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[s.Qualify(k)] = v
	}
	return out
}

// Unwrap returns a new map with the namespace prefix stripped from every key.
// Keys that do not carry the prefix are passed through unchanged.
func (s *Scoper) Unwrap(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[s.Strip(k)] = v
	}
	return out
}

// Prefix returns the namespace string used by this Scoper.
func (s *Scoper) Prefix() string { return s.prefix }
