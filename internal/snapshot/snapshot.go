// Package snapshot captures and restores secret maps for rollback support.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry holds a timestamped snapshot of secrets for a single path.
type Entry struct {
	Path      string            `json:"path"`
	Secrets   map[string]string `json:"secrets"`
	CreatedAt time.Time         `json:"created_at"`
}

// Store manages snapshot files on disk.
type Store struct {
	dir string
}

// NewStore creates a Store that persists snapshots under dir.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("snapshot: mkdir %s: %w", dir, err)
	}
	return &Store{dir: dir}, nil
}

// Save writes a snapshot for the given vault path and secrets.
func (s *Store) Save(vaultPath string, secrets map[string]string) error {
	e := Entry{
		Path:      vaultPath,
		Secrets:   secrets,
		CreatedAt: time.Now().UTC(),
	}
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	dest := s.filePath(vaultPath)
	if err := os.WriteFile(dest, data, 0600); err != nil {
		return fmt.Errorf("snapshot: write %s: %w", dest, err)
	}
	return nil
}

// Load returns the most recent snapshot for vaultPath, or nil if none exists.
func (s *Store) Load(vaultPath string) (*Entry, error) {
	data, err := os.ReadFile(s.filePath(vaultPath))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("snapshot: read: %w", err)
	}
	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &e, nil
}

func (s *Store) filePath(vaultPath string) string {
	safe := filepath.Clean(vaultPath)
	// replace slashes so each path maps to a flat filename
	for i := 0; i < len(safe); i++ {
		if safe[i] == '/' || safe[i] == os.PathSeparator {
			safe = safe[:i] + "_" + safe[i+1:]
		}
	}
	return filepath.Join(s.dir, safe+".json")
}
