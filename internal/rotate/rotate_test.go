package rotate

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRotate_NoSourceFile(t *testing.T) {
	dir := t.TempDir()
	r := New(filepath.Join(dir, "backups"), 3)
	if err := r.Rotate(filepath.Join(dir, "nonexistent.env")); err != nil {
		t.Fatalf("expected nil for missing source, got %v", err)
	}
}

func TestRotate_CreatesBackup(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, ".env")
	if err := os.WriteFile(src, []byte("KEY=val\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	backupDir := filepath.Join(dir, "backups")
	r := New(backupDir, 5)
	if err := r.Rotate(src); err != nil {
		t.Fatalf("Rotate: %v", err)
	}

	entries, _ := filepath.Glob(filepath.Join(backupDir, ".env.*.bak"))
	if len(entries) != 1 {
		t.Fatalf("expected 1 backup, got %d", len(entries))
	}

	data, _ := os.ReadFile(entries[0])
	if string(data) != "KEY=val\n" {
		t.Errorf("backup content mismatch: %q", data)
	}
}

func TestRotate_PrunesOldBackups(t *testing.T) {
	dir := t.TempDir()
	backupDir := filepath.Join(dir, "backups")
	if err := os.MkdirAll(backupDir, 0o700); err != nil {
		t.Fatal(err)
	}

	// Pre-create 5 old backups
	for i := 0; i < 5; i++ {
		ts := time.Now().UTC().Add(time.Duration(-5+i) * time.Minute).Format("20060102T150405Z")
		name := filepath.Join(backupDir, ".env."+ts+".bak")
		os.WriteFile(name, []byte("old"), 0o600)
		time.Sleep(time.Millisecond) // ensure unique names
	}

	src := filepath.Join(dir, ".env")
	os.WriteFile(src, []byte("NEW=1\n"), 0o600)

	r := New(backupDir, 3)
	if err := r.Rotate(src); err != nil {
		t.Fatalf("Rotate: %v", err)
	}

	entries, _ := filepath.Glob(filepath.Join(backupDir, ".env.*.bak"))
	if len(entries) != 3 {
		t.Errorf("expected 3 backups after prune, got %d", len(entries))
	}
}

func TestNew_DefaultMaxBackups(t *testing.T) {
	r := New("/tmp", 0)
	if r.maxBackups != 5 {
		t.Errorf("expected default 5, got %d", r.maxBackups)
	}
}
