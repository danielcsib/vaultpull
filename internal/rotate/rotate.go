package rotate

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Rotator manages backup rotation of .env files before overwriting.
type Rotator struct {
	backupDir string
	maxBackups int
}

// New creates a Rotator that stores backups in backupDir, keeping at most maxBackups.
func New(backupDir string, maxBackups int) *Rotator {
	if maxBackups <= 0 {
		maxBackups = 5
	}
	return &Rotator{backupDir: backupDir, maxBackups: maxBackups}
}

// Rotate copies src to a timestamped backup file, then prunes old backups.
func (r *Rotator) Rotate(src string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // nothing to back up
		}
		return fmt.Errorf("rotate: read %s: %w", src, err)
	}

	if err := os.MkdirAll(r.backupDir, 0o700); err != nil {
		return fmt.Errorf("rotate: mkdir %s: %w", r.backupDir, err)
	}

	base := filepath.Base(src)
	timestamp := time.Now().UTC().Format("20060102T150405Z")
	dest := filepath.Join(r.backupDir, fmt.Sprintf("%s.%s.bak", base, timestamp))

	if err := os.WriteFile(dest, data, 0o600); err != nil {
		return fmt.Errorf("rotate: write backup %s: %w", dest, err)
	}

	return r.prune(base)
}

// prune removes oldest backups for base beyond maxBackups.
func (r *Rotator) prune(base string) error {
	pattern := filepath.Join(r.backupDir, base+".*.bak")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	if len(matches) <= r.maxBackups {
		return nil
	}
	// matches are lexicographically sorted; oldest first
	for _, old := range matches[:len(matches)-r.maxBackups] {
		if err := os.Remove(old); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("rotate: prune %s: %w", old, err)
		}
	}
	return nil
}
