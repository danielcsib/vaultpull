// Package expire tracks TTL-based expiry for synced secret paths.
package expire

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry holds the expiry metadata for a single path.
type Entry struct {
	Path      string    `json:"path"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Store manages expiry entries persisted to disk.
type Store struct {
	mu      sync.Mutex
	file    string
	entries map[string]Entry
}

// NewStore loads or creates an expiry store at the given file path.
func NewStore(file string) (*Store, error) {
	s := &Store{file: file, entries: make(map[string]Entry)}
	data, err := os.ReadFile(file)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	var list []Entry
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	for _, e := range list {
		s.entries[e.Path] = e
	}
	return s, nil
}

// Set records an expiry for path with the given TTL from now.
func (s *Store) Set(path string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[path] = Entry{Path: path, ExpiresAt: time.Now().Add(ttl)}
	return s.save()
}

// Expired returns true if the entry for path exists and has expired.
func (s *Store) Expired(path string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[path]
	if !ok {
		return false
	}
	return time.Now().After(e.ExpiresAt)
}

// Delete removes the expiry entry for path.
func (s *Store) Delete(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, path)
	return s.save()
}

func (s *Store) save() error {
	list := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		list = append(list, e)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.file, data, 0600)
}
