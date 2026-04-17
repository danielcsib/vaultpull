package filter

import (
	"strings"
)

// Rule defines a single include or exclude rule for secret keys.
type Rule struct {
	Prefix  string
	Exclude bool
}

// Filter holds a set of rules to selectively include or exclude secret keys.
type Filter struct {
	rules []Rule
}

// New creates a Filter from a slice of rule strings.
// Strings prefixed with '!' are treated as exclusions.
func New(patterns []string) *Filter {
	rules := make([]Rule, 0, len(patterns))
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if strings.HasPrefix(p, "!") {
			rules = append(rules, Rule{Prefix: strings.TrimPrefix(p, "!"), Exclude: true})
		} else {
			rules = append(rules, Rule{Prefix: p, Exclude: false})
		}
	}
	return &Filter{rules: rules}
}

// Allow returns true if the given key should be included based on the rules.
// If no include rules exist, all keys are allowed unless excluded.
func (f *Filter) Allow(key string) bool {
	if len(f.rules) == 0 {
		return true
	}

	hasIncludes := false
	for _, r := range f.rules {
		if !r.Exclude {
			hasIncludes = true
			break
		}
	}

	for _, r := range f.rules {
		if r.Exclude && strings.HasPrefix(key, r.Prefix) {
			return false
		}
	}

	if !hasIncludes {
		return true
	}

	for _, r := range f.rules {
		if !r.Exclude && strings.HasPrefix(key, r.Prefix) {
			return true
		}
	}

	return false
}

// Apply filters a map of secrets, returning only allowed keys.
func (f *Filter) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if f.Allow(k) {
			out[k] = v
		}
	}
	return out
}
