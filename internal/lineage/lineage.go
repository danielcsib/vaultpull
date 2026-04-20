// Package lineage tracks the origin of each secret value — which Vault path
// and key it came from — so that downstream tooling can audit provenance.
package lineage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Entry records where a single env-key's value was sourced from.
type Entry struct {
	EnvKey    string    `json:"env_key"`
	VaultPath string    `json:"vault_path"`
	VaultKey  string    `json:"vault_key"`
	FetchedAt time.Time `json:"fetched_at"`
}

// Store holds lineage entries and can persist them to disk.
type Store struct {
	mu      sync.RWMutex
	entries map[string]Entry // keyed by EnvKey
	path    string
}

// NewStore creates a Store that persists to path.
// If path is empty the store operates in-memory only.
func NewStore(path string) (*Store, error) {
	s := &Store{path: path, entries: make(map[string]Entry)}
	if path == "" {
		return s, nil
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, fmt.Errorf("lineage: read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, &s.entries); err != nil {
		return nil, fmt.Errorf("lineage: parse %s: %w", path, err)
	}
	return s, nil
}

// Record stores an entry for the given env key.
func (s *Store) Record(envKey, vaultPath, vaultKey string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[envKey] = Entry{
		EnvKey:    envKey,
		VaultPath: vaultPath,
		VaultKey:  vaultKey,
		FetchedAt: time.Now().UTC(),
	}
}

// Get returns the lineage entry for envKey, if present.
func (s *Store) Get(envKey string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[envKey]
	return e, ok
}

// All returns a copy of all recorded entries.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, e)
	}
	return out
}

// Flush writes the current state to disk. No-op when path is empty.
func (s *Store) Flush() error {
	if s.path == "" {
		return nil
	}
	s.mu.RLock()
	data, err := json.MarshalIndent(s.entries, "", "  ")
	s.mu.RUnlock()
	if err != nil {
		return fmt.Errorf("lineage: marshal: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0o600); err != nil {
		return fmt.Errorf("lineage: write %s: %w", s.path, err)
	}
	return nil
}
