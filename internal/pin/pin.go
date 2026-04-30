// Package pin tracks pinned secret versions so that specific Vault
// secret versions are never silently overwritten during a sync.
package pin

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

// ErrAlreadyPinned is returned when a path is already pinned to a different version.
var ErrAlreadyPinned = errors.New("pin: path is already pinned to a different version")

// Entry holds a pinned version for a single secret path.
type Entry struct {
	Path      string    `json:"path"`
	Version   int       `json:"version"`
	PinnedAt  time.Time `json:"pinned_at"`
	PinnedBy  string    `json:"pinned_by"`
}

// Store persists pin entries to a JSON file.
type Store struct {
	mu   sync.RWMutex
	file string
	data map[string]Entry
}

// NewStore opens or creates a pin store backed by file.
func NewStore(file string) (*Store, error) {
	s := &Store{file: file, data: make(map[string]Entry)}
	if err := s.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("pin: load: %w", err)
	}
	return s, nil
}

// Pin records a pinned version for path. Returns ErrAlreadyPinned if the path
// is already pinned to a different version.
func (s *Store) Pin(path string, version int, by string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e, ok := s.data[path]; ok && e.Version != version {
		return fmt.Errorf("%w: %s@v%d", ErrAlreadyPinned, path, e.Version)
	}
	s.data[path] = Entry{Path: path, Version: version, PinnedAt: time.Now().UTC(), PinnedBy: by}
	return s.save()
}

// Unpin removes the pin for path. It is a no-op if the path is not pinned.
func (s *Store) Unpin(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, path)
	return s.save()
}

// Get returns the Entry for path and whether it exists.
func (s *Store) Get(path string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.data[path]
	return e, ok
}

// All returns a copy of all pinned entries.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.data))
	for _, e := range s.data {
		out = append(out, e)
	}
	return out
}

func (s *Store) load() error {
	b, err := os.ReadFile(s.file)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &s.data)
}

func (s *Store) save() error {
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.file, b, 0o600)
}
