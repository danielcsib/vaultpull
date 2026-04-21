// Package output formats and writes secret maps to various targets.
package output

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Format describes the output serialisation format.
type Format string

const (
	FormatEnv  Format = "env"
	FormatJSON Format = "json"
	FormatExport Format = "export"
)

// Writer serialises a secret map to a chosen format.
type Writer struct {
	format Format
	out    io.Writer
}

// New returns a Writer that writes to w using the given format.
// If w is nil, os.Stdout is used.
func New(format Format, w io.Writer) *Writer {
	if w == nil {
		w = os.Stdout
	}
	return &Writer{format: format, out: w}
}

// Write serialises secrets and writes them to the underlying writer.
func (w *Writer) Write(secrets map[string]string) error {
	switch w.format {
	case FormatEnv:
		return writeEnv(w.out, secrets, false)
	case FormatExport:
		return writeEnv(w.out, secrets, true)
	case FormatJSON:
		return writeJSON(w.out, secrets)
	default:
		return fmt.Errorf("output: unknown format %q", w.format)
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func writeEnv(w io.Writer, secrets map[string]string, export bool) error {
	prefix := ""
	if export {
		prefix = "export "
	}
	for _, k := range sortedKeys(secrets) {
		v := secrets[k]
		if strings.ContainsAny(v, " \t\n#") {
			v = fmt.Sprintf("%q", v)
		}
		if _, err := fmt.Fprintf(w, "%s%s=%s\n", prefix, k, v); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, secrets map[string]string) error {
	keys := sortedKeys(secrets)
	fmt.Fprintln(w, "{")
	for i, k := range keys {
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		if _, err := fmt.Fprintf(w, "  %q: %q%s\n", k, secrets[k], comma); err != nil {
			return err
		}
	}
	fmt.Fprintln(w, "}")
	return nil
}
