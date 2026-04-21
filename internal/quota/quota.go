// Package quota enforces per-path read limits to prevent runaway Vault requests.
package quota

import (
	"errors"
	"fmt"
	"sync"
)

// ErrQuotaExceeded is returned when a path has exceeded its allowed read count.
var ErrQuotaExceeded = errors.New("quota exceeded")

// Store tracks read counts per Vault secret path.
type Store struct {
	mu      sync.Mutex
	counts  map[string]int
	maxReads int
}

// New creates a Store that allows up to maxReads reads per path.
// If maxReads is zero or negative, it defaults to 10.
func New(maxReads int) *Store {
	if maxReads <= 0 {
		maxReads = 10
	}
	return &Store{
		counts:   make(map[string]int),
		maxReads: maxReads,
	}
}

// Check reports whether the given path is within quota.
// It increments the counter and returns ErrQuotaExceeded if the limit is hit.
func (s *Store) Check(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counts[path]++
	if s.counts[path] > s.maxReads {
		return fmt.Errorf("%w: path %q has been read %d times (max %d)",
			ErrQuotaExceeded, path, s.counts[path], s.maxReads)
	}
	return nil
}

// Remaining returns how many reads are left for the given path.
func (s *Store) Remaining(path string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	remaining := s.maxReads - s.counts[path]
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Reset clears the counter for a specific path.
func (s *Store) Reset(path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.counts, path)
}

// ResetAll clears all counters.
func (s *Store) ResetAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counts = make(map[string]int)
}

// Snapshot returns a copy of current read counts keyed by path.
func (s *Store) Snapshot() map[string]int {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make(map[string]int, len(s.counts))
	for k, v := range s.counts {
		out[k] = v
	}
	return out
}
