package rollback

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ListBackups prints all available backups for target to w.
func ListBackups(w io.Writer, backupDir, target string) error {
	if w == nil {
		w = os.Stdout
	}
	base := filepath.Base(target)
	pattern := filepath.Join(backupDir, base+".*.bak")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	if len(matches) == 0 {
		fmt.Fprintf(w, "no backups found for %s\n", target)
		return nil
	}
	fmt.Fprintf(w, "backups for %s (%d):\n", target, len(matches))
	for i, m := range matches {
		marker := "  "
		if i == len(matches)-1 {
			marker = "* "
		}
		name := strings.TrimPrefix(filepath.Base(m), base+".")
		name = strings.TrimSuffix(name, ".bak")
		fmt.Fprintf(w, "%s%s\n", marker, name)
	}
	return nil
}
