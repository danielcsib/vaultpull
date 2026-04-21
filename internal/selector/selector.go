// Package selector provides key-based selection over secret maps,
// allowing callers to pick a subset of keys from a secrets map.
package selector

import (
	"fmt"
	"sort"
)

// Selector holds a set of allowed keys for filtering secret maps.
type Selector struct {
	keys map[string]struct{}
}

// New creates a Selector from the provided list of keys.
// Duplicate keys are silently deduplicated.
// Returns an error if no keys are provided.
func New(keys []string) (*Selector, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("selector: at least one key must be specified")
	}
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		if k == "" {
			return nil, fmt.Errorf("selector: empty key is not allowed")
		}
		m[k] = struct{}{}
	}
	return &Selector{keys: m}, nil
}

// Pick returns a new map containing only the entries whose keys are
// present in the selector. Keys that are selected but absent from
// secrets are recorded in the returned missing slice.
func (s *Selector) Pick(secrets map[string]string) (picked map[string]string, missing []string) {
	picked = make(map[string]string)
	for k := range s.keys {
		v, ok := secrets[k]
		if !ok {
			missing = append(missing, k)
			continue
		}
		picked[k] = v
	}
	sort.Strings(missing)
	return picked, missing
}

// Has reports whether the given key is in the selector.
func (s *Selector) Has(key string) bool {
	_, ok := s.keys[key]
	return ok
}

// Keys returns the sorted list of keys held by the selector.
func (s *Selector) Keys() []string {
	out := make([]string, 0, len(s.keys))
	for k := range s.keys {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
