// Package flatten provides utilities for flattening nested secret maps
// into a single-level key=value map suitable for .env file output.
package flatten

import (
	"fmt"
	"strings"
)

// Options controls how keys are constructed during flattening.
type Options struct {
	// Separator is placed between nested key segments. Defaults to "_".
	Separator string
	// UppercaseKeys converts all keys to uppercase when true.
	UppercaseKeys bool
}

// DefaultOptions returns sensible defaults for flattening.
func DefaultOptions() Options {
	return Options{
		Separator:     "_",
		UppercaseKeys: true,
	}
}

// Flattener collapses nested map[string]any structures into a flat
// map[string]string using a configurable separator.
type Flattener struct {
	opts Options
}

// New creates a Flattener with the provided options.
func New(opts Options) *Flattener {
	if opts.Separator == "" {
		opts.Separator = "_"
	}
	return &Flattener{opts: opts}
}

// Flatten takes a nested map and returns a flat map[string]string.
// Nested maps are recursively expanded; all other values are converted
// to strings via fmt.Sprintf.
func (f *Flattener) Flatten(input map[string]any) map[string]string {
	out := make(map[string]string)
	f.flatten("", input, out)
	return out
}

func (f *Flattener) flatten(prefix string, input map[string]any, out map[string]string) {
	for k, v := range input {
		key := f.buildKey(prefix, k)
		switch val := v.(type) {
		case map[string]any:
			f.flatten(key, val, out)
		case map[string]string:
			for sk, sv := range val {
				out[f.buildKey(key, sk)] = sv
			}
		default:
			out[key] = fmt.Sprintf("%v", val)
		}
	}
}

func (f *Flattener) buildKey(prefix, segment string) string {
	var key string
	if prefix == "" {
		key = segment
	} else {
		key = prefix + f.opts.Separator + segment
	}
	if f.opts.UppercaseKeys {
		return strings.ToUpper(key)
	}
	return key
}
