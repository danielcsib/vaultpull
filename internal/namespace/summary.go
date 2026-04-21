package namespace

import (
	"fmt"
	"io"
	"sort"
)

// PrintSummary writes a human-readable table of namespaced keys to w.
// Each row shows the original key and its qualified form.
func PrintSummary(w io.Writer, s *Scoper, keys []string) {
	if len(keys) == 0 {
		fmt.Fprintln(w, "namespace: no keys to display")
		return
	}

	sorted := make([]string, len(keys))
	copy(sorted, keys)
	sort.Strings(sorted)

	fmt.Fprintf(w, "%-30s  %s\n", "KEY", "QUALIFIED")
	fmt.Fprintf(w, "%-30s  %s\n", "---", "---------")
	for _, k := range sorted {
		fmt.Fprintf(w, "%-30s  %s\n", k, s.Qualify(k))
	}
}
