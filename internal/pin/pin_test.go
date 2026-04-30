package pin_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/pin"
)

func tempFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "pins.json")
}

func TestNewStore_MissingFile(t *testing.T) {
	_, err := pin.NewStore(tempFile(t))
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
}

func TestPin_And_Get(t *testing.T) {
	s, _ := pin.NewStore(tempFile(t))
	if err := s.Pin("secret/app", 3, "alice"); err != nil {
		t.Fatalf("Pin: %v", err)
	}
	e, ok := s.Get("secret/app")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Version != 3 || e.PinnedBy != "alice" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestPin_SameVersion_IsIdempotent(t *testing.T) {
	s, _ := pin.NewStore(tempFile(t))
	s.Pin("secret/app", 2, "alice")
	if err := s.Pin("secret/app", 2, "bob"); err != nil {
		t.Errorf("expected no error for same version re-pin, got %v", err)
	}
}

func TestPin_DifferentVersion_ReturnsError(t *testing.T) {
	s, _ := pin.NewStore(tempFile(t))
	s.Pin("secret/app", 1, "alice")
	err := s.Pin("secret/app", 2, "bob")
	if err == nil {
		t.Fatal("expected ErrAlreadyPinned")
	}
}

func TestUnpin_RemovesEntry(t *testing.T) {
	s, _ := pin.NewStore(tempFile(t))
	s.Pin("secret/app", 1, "alice")
	s.Unpin("secret/app")
	_, ok := s.Get("secret/app")
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	s, _ := pin.NewStore(tempFile(t))
	s.Pin("secret/a", 1, "alice")
	s.Pin("secret/b", 2, "bob")
	if got := len(s.All()); got != 2 {
		t.Errorf("expected 2 entries, got %d", got)
	}
}

func TestStore_Persists(t *testing.T) {
	f := tempFile(t)
	s1, _ := pin.NewStore(f)
	s1.Pin("secret/db", 5, "ci")

	s2, err := pin.NewStore(f)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	e, ok := s2.Get("secret/db")
	if !ok || e.Version != 5 {
		t.Errorf("expected persisted entry, got %+v", e)
	}
}

func TestNewStore_CorruptFile(t *testing.T) {
	f := tempFile(t)
	os.WriteFile(f, []byte("not-json{"), 0o600)
	_, err := pin.NewStore(f)
	if err == nil {
		t.Error("expected error for corrupt file")
	}
}
