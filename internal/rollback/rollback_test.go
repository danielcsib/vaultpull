package rollback_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/rollback"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "rollback-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestLatest_NoBackups(t *testing.T) {
	dir := tempDir(t)
	rb := rollback.New(dir)
	_, err := rb.Latest("/some/.env")
	if err == nil {
		t.Fatal("expected error for missing backups")
	}
}

func TestLatest_ReturnsLatest(t *testing.T) {
	dir := tempDir(t)
	writeFile(t, filepath.Join(dir, ".env.20240101T120000.bak"), "A=1")
	writeFile(t, filepath.Join(dir, ".env.20240102T120000.bak"), "A=2")
	rb := rollback.New(dir)
	got, err := rb.Latest(".env")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(got) != ".env.20240102T120000.bak" {
		t.Errorf("unexpected latest: %s", got)
	}
}

func TestRestore_WritesContent(t *testing.T) {
	dir := tempDir(t)
	backup := filepath.Join(dir, ".env.20240101T120000.bak")
	writeFile(t, backup, "RESTORED=true")
	target := filepath.Join(dir, ".env")
	writeFile(t, target, "OLD=value")
	rb := rollback.New(dir)
	_, err := rb.Restore(target)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(target)
	if string(data) != "RESTORED=true" {
		t.Errorf("unexpected content: %s", data)
	}
}

func TestRestoreFrom_CopiesFile(t *testing.T) {
	dir := tempDir(t)
	src := filepath.Join(dir, "backup.bak")
	dst := filepath.Join(dir, ".env")
	writeFile(t, src, "KEY=hello")
	rb := rollback.New(dir)
	if err := rb.RestoreFrom(src, dst); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(dst)
	if string(data) != "KEY=hello" {
		t.Errorf("got %s", data)
	}
}
