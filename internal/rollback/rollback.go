package rollback

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Rollbacker restores a .env file from a backup.
type Rollbacker struct {
	backupDir string
}

// New returns a Rollbacker that reads backups from backupDir.
func New(backupDir string) *Rollbacker {
	return &Rollbacker{backupDir: backupDir}
}

// Latest returns the path of the most recently modified backup for target.
func (r *Rollbacker) Latest(target string) (string, error) {
	base := filepath.Base(target)
	pattern := filepath.Join(r.backupDir, base+".*.bak")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("rollback: no backups found for %s", target)
	}
	// Glob returns sorted; last entry is lexicographically latest (timestamp suffix).
	return matches[len(matches)-1], nil
}

// Restore copies the latest backup over target.
func (r *Rollbacker) Restore(target string) (string, error) {
	src, err := r.Latest(target)
	if err != nil {
		return "", err
	}
	if err := copyFile(src, target); err != nil {
		return "", fmt.Errorf("rollback: restore failed: %w", err)
	}
	return src, nil
}

// RestoreFrom copies a specific backup file over target.
func (r *Rollbacker) RestoreFrom(src, target string) error {
	return copyFile(src, target)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
