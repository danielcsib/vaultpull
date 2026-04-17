package redact

import "strings"

// Redactor masks sensitive secret values before logging or displaying.
type Redactor struct {
	visibleChars int
}

// New returns a Redactor that shows at most visibleChars trailing characters.
func New(visibleChars int) *Redactor {
	if visibleChars < 0 {
		visibleChars = 0
	}
	return &Redactor{visibleChars: visibleChars}
}

// Mask replaces most of value with asterisks, revealing only the last
// visibleChars characters. Empty strings are returned as-is.
func (r *Redactor) Mask(value string) string {
	if value == "" {
		return ""
	}
	if r.visibleChars == 0 || len(value) <= r.visibleChars {
		return strings.Repeat("*", len(value))
	}
	hidden := len(value) - r.visibleChars
	return strings.Repeat("*", hidden) + value[hidden:]
}

// MaskMap returns a new map with all values masked.
func (r *Redactor) MaskMap(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = r.Mask(v)
	}
	return out
}
