package expire_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/expire"
)

func tempFile(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "expire.json")
}

func TestNewStore_MissingFile(t *testing.T) {
	s, err := expire.NewStore(tempFile(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil store")
	}
}

func TestSetAndExpired_NotYetExpired(t *testing.T) {
	s, _ := expire.NewStore(tempFile(t))
	if err := s.Set("secret/app", 10*time.Minute); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if s.Expired("secret/app") {
		t.Error("expected not expired")
	}
}

func TestExpired_AlreadyExpired(t *testing.T) {
	s, _ := expire.NewStore(tempFile(t))
	_ = s.Set("secret/old", -1*time.Second)
	if !s.Expired("secret/old") {
		t.Error("expected expired")
	}
}

func TestExpired_UnknownPath(t *testing.T) {
	s, _ := expire.NewStore(tempFile(t))
	if s.Expired("secret/unknown") {
		t.Error("unknown path should not be expired")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	s, _ := expire.NewStore(tempFile(t))
	_ = s.Set("secret/tmp", -1*time.Second)
	_ = s.Delete("secret/tmp")
	if s.Expired("secret/tmp") {
		t.Error("deleted entry should not be expired")
	}
}

func TestPersistence_RoundTrip(t *testing.T) {
	f := tempFile(t)
	s1, _ := expire.NewStore(f)
	_ = s1.Set("secret/db", 1*time.Hour)

	s2, err := expire.NewStore(f)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if s2.Expired("secret/db") {
		t.Error("reloaded entry should not be expired")
	}
}

func TestNewStore_CorruptFile(t *testing.T) {
	f := tempFile(t)
	_ = os.WriteFile(f, []byte("not json"), 0600)
	_, err := expire.NewStore(f)
	if err == nil {
		t.Error("expected error for corrupt file")
	}
}
