// Package sanitize provides utilities for cleaning and normalising secret
// values before they are written to .env files.
package sanitize

import (
	"strings"
	"unicode"
)

// Options controls which sanitisation steps are applied.
type Options struct {
	// TrimSpace removes leading and trailing whitespace from values.
	TrimSpace bool
	// StripNonPrintable removes non-printable Unicode characters from values.
	StripNonPrintable bool
	// NormaliseNewlines replaces CR+LF and bare CR sequences with LF.
	NormaliseNewlines bool
}

// DefaultOptions returns a sensible default set of sanitisation options.
func DefaultOptions() Options {
	return Options{
		TrimSpace:         true,
		StripNonPrintable: true,
		NormaliseNewlines: true,
	}
}

// Sanitizer applies a fixed set of Options to secret values.
type Sanitizer struct {
	opts Options
}

// New returns a Sanitizer configured with opts.
func New(opts Options) *Sanitizer {
	return &Sanitizer{opts: opts}
}

// Value sanitises a single string value according to the configured Options.
func (s *Sanitizer) Value(v string) string {
	if s.opts.NormaliseNewlines {
		v = strings.ReplaceAll(v, "\r\n", "\n")
		v = strings.ReplaceAll(v, "\r", "\n")
	}
	if s.opts.StripNonPrintable {
		v = stripNonPrintable(v)
	}
	if s.opts.TrimSpace {
		v = strings.TrimSpace(v)
	}
	return v
}

// Apply sanitises every value in m, returning a new map. The original map is
// never mutated.
func (s *Sanitizer) Apply(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = s.Value(v)
	}
	return out
}

func stripNonPrintable(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r == '\n' || r == '\t' || unicode.IsPrint(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}
