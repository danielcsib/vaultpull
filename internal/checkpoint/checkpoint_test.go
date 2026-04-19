package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/checkpoint"
)

func tempFile(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "checkpoint.json")
}

func TestGet_MissingEntry(t *testing.T) {
	s, err := checkpoint.NewStore(tempFile(t))
	if err != nil {
		t.Fatal(err)
	}
	_, ok := s.Get("secret/app")
	if ok {
		t.Fatal("expected no entry")
	}
}

func TestRecord_And_Get(t *testing.T) {
	s, err := checkpoint.NewStore(tempFile(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Record("secret/app", 5); err != nil {
		t.Fatal(err)
	}
	e, ok := s.Get("secret/app")
	if !ok {
		t.Fatal("expected entry")
	}
	if e.KeyCount != 5 {
		t.Errorf("got %d keys, want 5", e.KeyCount)
	}
	if e.SyncedAt.IsZero() {
		t.Error("expected non-zero SyncedAt")
	}
}

func TestRecord_Persists(t *testing.T) {
	path := tempFile(t)
	s, _ := checkpoint.NewStore(path)
	s.Record("secret/db", 3)

	s2, err := checkpoint.NewStore(path)
	if err != nil {
		t.Fatal(err)
	}
	e, ok := s2.Get("secret/db")
	if !ok {
		t.Fatal("entry not persisted")
	}
	if e.KeyCount != 3 {
		t.Errorf("got %d, want 3", e.KeyCount)
	}
}

func TestAll_ReturnsEntries(t *testing.T) {
	s, _ := checkpoint.NewStore(tempFile(t))
	s.Record("secret/a", 1)
	s.Record("secret/b", 2)
	all := s.All()
	if len(all) != 2 {
		t.Errorf("got %d entries, want 2", len(all))
	}
}

func TestNewStore_CreatesDirectory(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "sub", "dir", "cp.json")
	_, err := checkpoint.NewStore(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		t.Error("directory not created")
	}
}

func TestRecord_UpdatesExisting(t *testing.T) {
	s, _ := checkpoint.NewStore(tempFile(t))
	s.Record("secret/x", 2)
	time.Sleep(time.Millisecond)
	s.Record("secret/x", 7)
	e, _ := s.Get("secret/x")
	if e.KeyCount != 7 {
		t.Errorf("got %d, want 7", e.KeyCount)
	}
}
