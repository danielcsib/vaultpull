package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/snapshot"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "snapshot-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestLoad_NoExistingSnapshot(t *testing.T) {
	st, _ := snapshot.NewStore(tempDir(t))
	e, err := st.Load("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e != nil {
		t.Fatal("expected nil entry for missing snapshot")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	st, _ := snapshot.NewStore(tempDir(t))
	secrets := map[string]string{"DB_PASS": "s3cr3t", "API_KEY": "abc123"}

	if err := st.Save("secret/app", secrets); err != nil {
		t.Fatalf("Save: %v", err)
	}
	e, err := st.Load("secret/app")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if e == nil {
		t.Fatal("expected entry, got nil")
	}
	if e.Path != "secret/app" {
		t.Errorf("path = %q, want %q", e.Path, "secret/app")
	}
	if e.Secrets["DB_PASS"] != "s3cr3t" {
		t.Errorf("DB_PASS = %q, want s3cr3t", e.Secrets["DB_PASS"])
	}
	if e.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestSave_OverwritesPrevious(t *testing.T) {
	st, _ := snapshot.NewStore(tempDir(t))
	_ = st.Save("secret/app", map[string]string{"K": "v1"})
	_ = st.Save("secret/app", map[string]string{"K": "v2"})
	e, _ := st.Load("secret/app")
	if e.Secrets["K"] != "v2" {
		t.Errorf("expected v2, got %q", e.Secrets["K"])
	}
}

func TestNewStore_CreatesDirectory(t *testing.T) {
	base := tempDir(t)
	nested := filepath.Join(base, "a", "b", "snapshots")
	_, err := snapshot.NewStore(nested)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	if _, err := os.Stat(nested); os.IsNotExist(err) {
		t.Error("directory was not created")
	}
}
