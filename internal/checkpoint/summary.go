package checkpoint

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// PrintSummary writes a formatted table of all checkpoint entries to w.
func PrintSummary(w io.Writer, entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "no checkpoints recorded")
		return
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PATH\tKEYS\tLAST SYNCED")
	fmt.Fprintln(tw, "----\t-----\t-----------")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%d\t%s\n",
			e.Path,
			e.KeyCount,
			e.SyncedAt.Format("2006-01-02 15:04:05 UTC"),
		)
	}
	tw.Flush()
}
