package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a cached secret snapshot.
type Entry struct {
	Hash      string            `json:"hash"`
	FetchedAt time.Time         `json:"fetched_at"`
	Secrets   map[string]string `json:"secrets"`
}

// Store manages on-disk secret caches.
type Store struct {
	dir string
}

// NewStore creates a Store that persists cache files under dir.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("cache: create dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// key returns a filesystem-safe cache key for a vault path.
func (s *Store) key(vaultPath string) string {
	h := sha256.Sum256([]byte(vaultPath))
	return hex.EncodeToString(h[:]) + ".json"
}

// Get retrieves a cached entry for vaultPath. Returns nil if not found.
func (s *Store) Get(vaultPath string) (*Entry, error) {
	p := filepath.Join(s.dir, s.key(vaultPath))
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cache: read: %w", err)
	}
	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, fmt.Errorf("cache: unmarshal: %w", err)
	}
	return &e, nil
}

// Set writes secrets for vaultPath into the cache.
func (s *Store) Set(vaultPath string, secrets map[string]string) error {
	e := Entry{
		Hash:      hashSecrets(secrets),
		FetchedAt: time.Now().UTC(),
		Secrets:   secrets,
	}
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Errorf("cache: marshal: %w", err)
	}
	p := filepath.Join(s.dir, s.key(vaultPath))
	return os.WriteFile(p, data, 0600)
}

// Changed reports whether secrets differ from the cached entry.
func (s *Store) Changed(vaultPath string, secrets map[string]string) (bool, error) {
	e, err := s.Get(vaultPath)
	if err != nil {
		return true, err
	}
	if e == nil {
		return true, nil
	}
	return e.Hash != hashSecrets(secrets), nil
}

func hashSecrets(secrets map[string]string) string {
	b, _ := json.Marshal(secrets)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
