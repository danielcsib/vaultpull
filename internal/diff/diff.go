package diff

import "sort"

// Result holds the outcome of comparing two secret maps.
type Result struct {
	Added   []string
	Removed []string
	Changed []string
	Unchanged []string
}

// Compare returns a Result describing the difference between old and new
// secret maps. Keys map to their string values.
func Compare(old, next map[string]string) Result {
	r := Result{}

	for k, nv := range next {
		ov, exists := old[k]
		switch {
		case !exists:
			r.Added = append(r.Added, k)
		case ov != nv:
			r.Changed = append(r.Changed, k)
		default:
			r.Unchanged = append(r.Unchanged, k)
		}
	}

	for k := range old {
		if _, exists := next[k]; !exists {
			r.Removed = append(r.Removed, k)
		}
	}

	sort.Strings(r.Added)
	sort.Strings(r.Removed)
	sort.Strings(r.Changed)
	sort.Strings(r.Unchanged)

	return r
}

// HasChanges returns true when any keys were added, removed, or changed.
func (r Result) HasChanges() bool {
	return len(r.Added) > 0 || len(r.Removed) > 0 || len(r.Changed) > 0
}
