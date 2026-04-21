// Package tag provides key tagging and lookup for annotating secrets
// with arbitrary metadata such as environment, team, or sensitivity level.
package tag

import "fmt"

// Tag represents a single key=value annotation.
type Tag struct {
	Key   string
	Value string
}

// Tagger holds a registry of tags per secret key.
type Tagger struct {
	entries map[string][]Tag
}

// New returns an empty Tagger.
func New() *Tagger {
	return &Tagger{entries: make(map[string][]Tag)}
}

// Set attaches a tag to a secret key. Duplicate tag keys for the same
// secret key are overwritten.
func (t *Tagger) Set(secretKey, tagKey, tagValue string) {
	if secretKey == "" || tagKey == "" {
		return
	}
	existing := t.entries[secretKey]
	for i, tag := range existing {
		if tag.Key == tagKey {
			existing[i].Value = tagValue
			t.entries[secretKey] = existing
			return
		}
	}
	t.entries[secretKey] = append(existing, Tag{Key: tagKey, Value: tagValue})
}

// Get returns all tags attached to a secret key.
func (t *Tagger) Get(secretKey string) []Tag {
	return t.entries[secretKey]
}

// Has reports whether a secret key carries a tag with the given key and value.
func (t *Tagger) Has(secretKey, tagKey, tagValue string) bool {
	for _, tag := range t.entries[secretKey] {
		if tag.Key == tagKey && tag.Value == tagValue {
			return true
		}
	}
	return false
}

// Filter returns the subset of keys from m whose tags satisfy all provided
// tag constraints (AND semantics).
func (t *Tagger) Filter(m map[string]string, constraints []Tag) map[string]string {
	out := make(map[string]string)
	for k, v := range m {
		if t.matches(k, constraints) {
			out[k] = v
		}
	}
	return out
}

func (t *Tagger) matches(secretKey string, constraints []Tag) bool {
	for _, c := range constraints {
		if !t.Has(secretKey, c.Key, c.Value) {
			return false
		}
	}
	return true
}

// String returns a human-readable representation of all tags for a key.
func (t *Tagger) String(secretKey string) string {
	tags := t.Get(secretKey)
	if len(tags) == 0 {
		return fmt.Sprintf("%s: (no tags)", secretKey)
	}
	s := secretKey + ":"
	for _, tag := range tags {
		s += fmt.Sprintf(" %s=%s", tag.Key, tag.Value)
	}
	return s
}
