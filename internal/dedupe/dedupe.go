// Package dedupe provides deduplication of secret maps by tracking seen keys
// and resolving conflicts between multiple Vault sources.
package dedupe

import "fmt"

// ConflictPolicy controls how key conflicts are resolved.
type ConflictPolicy int

const (
	// PolicyKeepFirst retains the first value seen for a key.
	PolicyKeepFirst ConflictPolicy = iota
	// PolicyKeepLast overwrites with the latest value seen.
	PolicyKeepLast
	// PolicyError returns an error on any conflict.
	PolicyError
)

// Conflict records a key that appeared in more than one source.
type Conflict struct {
	Key    string
	First  string
	Second string
}

// Merger merges multiple secret maps according to a conflict policy.
type Merger struct {
	policy    ConflictPolicy
	Conflicts []Conflict
}

// New returns a Merger with the given policy.
func New(policy ConflictPolicy) *Merger {
	return &Merger{policy: policy}
}

// Merge combines src into dst, applying the conflict policy when keys collide.
func (m *Merger) Merge(dst, src map[string]string) error {
	for k, v := range src {
		existing, exists := dst[k]
		if !exists {
			dst[k] = v
			continue
		}
		m.Conflicts = append(m.Conflicts, Conflict{Key: k, First: existing, Second: v})
		switch m.policy {
		case PolicyKeepFirst:
			// retain existing — do nothing
		case PolicyKeepLast:
			dst[k] = v
		case PolicyError:
			return fmt.Errorf("dedupe: conflict on key %q", k)
		}
	}
	return nil
}

// Reset clears recorded conflicts.
func (m *Merger) Reset() {
	m.Conflicts = nil
}
