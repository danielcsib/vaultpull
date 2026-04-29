// Package coalesce provides utilities for merging multiple secret maps
// with a defined priority order, returning the first non-empty value for
// each key across the provided sources.
package coalesce

import "errors"

// ErrNoSources is returned when Merge is called with no source maps.
var ErrNoSources = errors.New("coalesce: at least one source map is required")

// Merger merges multiple secret maps in priority order.
type Merger struct {
	omitEmpty bool
}

// Option configures a Merger.
type Option func(*Merger)

// OmitEmpty causes keys with empty string values to be skipped entirely
// in the merged output.
func OmitEmpty() Option {
	return func(m *Merger) { m.omitEmpty = true }
}

// New creates a new Merger with the given options.
func New(opts ...Option) *Merger {
	m := &Merger{}
	for _, o := range opts {
		o(m)
	}
	return m
}

// Merge combines sources in priority order: sources[0] has highest priority.
// For each key, the first non-empty value found across sources is used.
// If omitEmpty is false, an empty string value from a higher-priority source
// still wins over a non-empty value in a lower-priority source.
func (m *Merger) Merge(sources ...map[string]string) (map[string]string, error) {
	if len(sources) == 0 {
		return nil, ErrNoSources
	}

	out := make(map[string]string)

	// Collect all keys across all sources.
	keys := make(map[string]struct{})
	for _, src := range sources {
		for k := range src {
			keys[k] = struct{}{}
		}
	}

	for key := range keys {
		for _, src := range sources {
			v, ok := src[key]
			if !ok {
				continue
			}
			if v == "" && m.omitEmpty {
				continue
			}
			out[key] = v
			break
		}
	}

	return out, nil
}

// First returns the first map in sources that contains the given key,
// along with its value. Returns ("", false) if no source contains the key.
func First(key string, sources ...map[string]string) (string, bool) {
	for _, src := range sources {
		if v, ok := src[key]; ok {
			return v, true
		}
	}
	return "", false
}
