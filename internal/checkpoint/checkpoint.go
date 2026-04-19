// Package checkpoint tracks which vault paths were last synced successfully.
package checkpoint

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Entry records the last successful sync for a vault path.
type Entry struct {
	Path      string    `json:"path"`
	SyncedAt  time.Time `json:"synced_at"`
	KeyCount  int       `json:"key_count"`
}

// Store persists checkpoint entries to disk.
type Store struct {
	filePath string
	entries  map[string]Entry
}

// NewStore opens or creates a checkpoint store at the given file path.
func NewStore(filePath string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(filePath), 0700); err != nil {
		return nil, err
	}
	s := &Store{filePath: filePath, entries: make(map[string]Entry)}
	data, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &s.entries); err != nil {
		return nil, err
	}
	return s, nil
}

// Record saves a successful sync entry for the given path.
func (s *Store) Record(path string, keyCount int) error {
	s.entries[path] = Entry{
		Path:     path,
		SyncedAt: time.Now().UTC(),
		KeyCount: keyCount,
	}
	return s.flush()
}

// Get returns the checkpoint entry for a path, if any.
func (s *Store) Get(path string) (Entry, bool) {
	e, ok := s.entries[path]
	return e, ok
}

// All returns all recorded entries.
func (s *Store) All() []Entry {
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, e)
	}
	return out
}

func (s *Store) flush() error {
	data, err := json.MarshalIndent(s.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0600)
}
