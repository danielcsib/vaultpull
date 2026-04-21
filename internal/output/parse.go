package output

import (
	"fmt"
	"strings"
)

// ParseFormat converts a raw string into a Format constant.
// It returns an error if the value is not recognised.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "env", "":
		return FormatEnv, nil
	case "export":
		return FormatExport, nil
	case "json":
		return FormatJSON, nil
	default:
		return "", fmt.Errorf("output: unrecognised format %q; valid choices are env, export, json", s)
	}
}

// MustParseFormat is like ParseFormat but panics on an invalid value.
// Intended for use in tests and flag defaults.
func MustParseFormat(s string) Format {
	f, err := ParseFormat(s)
	if err != nil {
		panic(err)
	}
	return f
}
