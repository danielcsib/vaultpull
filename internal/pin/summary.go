package pin

import (
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"
)

// PrintSummary writes a human-readable table of pinned secrets to w.
// If w is nil, os.Stdout is used.
func PrintSummary(s *Store, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	entries := s.All()
	if len(entries) == 0 {
		fmt.Fprintln(w, "no pinned secrets")
		return
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PATH\tVERSION\tPINNED BY\tPINNED AT")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\tv%d\t%s\t%s\n",
			e.Path, e.Version, e.PinnedBy, e.PinnedAt.Format("2006-01-02 15:04:05 UTC"))
	}
	tw.Flush()
}
