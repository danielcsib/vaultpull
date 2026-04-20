package lineage_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/lineage"
)

func tempFile(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "lineage.json")
}

func TestNewStore_InMemory(t *testing.T) {
	s, err := lineage.NewStore("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil store")
	}
}

func TestRecord_And_Get(t *testing.T) {
	s, _ := lineage.NewStore("")
	s.Record("DB_PASSWORD", "secret/db", "password")

	e, ok := s.Get("DB_PASSWORD")
	if !ok {
		t.Fatal("expected entry to be present")
	}
	if e.VaultPath != "secret/db" {
		t.Errorf("VaultPath: got %q, want %q", e.VaultPath, "secret/db")
	}
	if e.VaultKey != "password" {
		t.Errorf("VaultKey: got %q, want %q", e.VaultKey, "password")
	}
	if e.FetchedAt.IsZero() {
		t.Error("FetchedAt should not be zero")
	}
}

func TestGet_MissingKey(t *testing.T) {
	s, _ := lineage.NewStore("")
	_, ok := s.Get("NONEXISTENT")
	if ok {
		t.Error("expected false for missing key")
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	s, _ := lineage.NewStore("")
	s.Record("KEY_A", "secret/a", "val")
	s.Record("KEY_B", "secret/b", "val")

	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestFlush_And_Reload(t *testing.T) {
	p := tempFile(t)
	s, _ := lineage.NewStore(p)
	s.Record("API_KEY", "secret/api", "key")

	if err := s.Flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}

	s2, err := lineage.NewStore(p)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	e, ok := s2.Get("API_KEY")
	if !ok {
		t.Fatal("expected entry after reload")
	}
	if e.VaultPath != "secret/api" {
		t.Errorf("VaultPath after reload: got %q", e.VaultPath)
	}
}

func TestFlush_NoOp_WhenPathEmpty(t *testing.T) {
	s, _ := lineage.NewStore("")
	s.Record("X", "p", "k")
	if err := s.Flush(); err != nil {
		t.Errorf("expected no error for in-memory flush, got %v", err)
	}
}

func TestNewStore_CorruptFile(t *testing.T) {
	p := tempFile(t)
	if err := os.WriteFile(p, []byte("not-json{"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := lineage.NewStore(p)
	if err == nil {
		t.Error("expected error for corrupt file")
	}
}
