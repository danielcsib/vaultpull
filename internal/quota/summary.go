package quota

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// PrintSummary writes a human-readable quota usage table to w.
// If w is nil, os.Stdout is used.
func PrintSummary(s *Store, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}

	snap := s.Snapshot()
	if len(snap) == 0 {
		fmt.Fprintln(w, "quota: no paths have been read")
		return
	}

	paths := make([]string, 0, len(snap))
	for p := range snap {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	fmt.Fprintln(w, "quota usage:")
	for _, p := range paths {
		count := snap[p]
		remaining := s.Remaining(p)
		status := "ok"
		if remaining == 0 {
			status = "EXCEEDED"
		}
		fmt.Fprintf(w, "  %-40s reads=%-4d remaining=%-4d [%s]\n", p, count, remaining, status)
	}
}
