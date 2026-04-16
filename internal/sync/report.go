package sync

import (
	"fmt"
	"io"
)

// PrintReport writes a human-readable summary of sync results to w.
func PrintReport(w io.Writer, results []Result) {
	successCount := 0
	failCount := 0

	for _, r := range results {
		if r.Err != nil {
			fmt.Fprintf(w, "  ✗ %s → %s: %v\n", r.Mapping.VaultPath, r.Mapping.EnvFile, r.Err)
			failCount++
		} else {
			fmt.Fprintf(w, "  ✓ %s → %s (%d keys)\n", r.Mapping.VaultPath, r.Mapping.EnvFile, r.Written)
			successCount++
		}
	}

	fmt.Fprintf(w, "\nDone: %d succeeded, %d failed.\n", successCount, failCount)
}

// HasErrors returns true if any result contains an error.
func HasErrors(results []Result) bool {
	for _, r := range results {
		if r.Err != nil {
			return true
		}
	}
	return false
}
