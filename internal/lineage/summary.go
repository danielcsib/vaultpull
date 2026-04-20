package lineage

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// PrintSummary writes a human-readable table of lineage entries to w.
func PrintSummary(w io.Writer, s *Store) {
	entries := s.All()
	if len(entries) == 0 {
		fmt.Fprintln(w, "no lineage entries recorded")
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].EnvKey < entries[j].EnvKey
	})

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ENV KEY\tVAULT PATH\tVAULT KEY\tFETCHED AT")
	fmt.Fprintln(tw, "-------\t----------\t---------\t----------")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			e.EnvKey,
			e.VaultPath,
			e.VaultKey,
			e.FetchedAt.Format("2006-01-02T15:04:05Z"),
		)
	}
	tw.Flush()
}
