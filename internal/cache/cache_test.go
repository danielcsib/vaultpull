package cache

import (
	"os"
	"testing"
)

func tempStore(t *testing.T) *Store {
	t.Helper()
	dir, err := os.MkdirTemp("", "vaultpull-cache-*")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	s, err := NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func TestGet_MissingEntry(t *testing.T) {
	s := tempStore(t)
	e, err := s.Get("secret/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e != nil {
		t.Fatal("expected nil entry for missing path")
	}
}

func TestSetAndGet(t *testing.T) {
	s := tempStore(t)
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := s.Set("secret/myapp", secrets); err != nil {
		t.Fatalf("Set: %v", err)
	}
	e, err := s.Get("secret/myapp")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if e == nil {
		t.Fatal("expected entry, got nil")
	}
	if e.Secrets["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", e.Secrets["FOO"])
	}
}

func TestChanged_NewPath(t *testing.T) {
	s := tempStore(t)
	changed, err := s.Changed("secret/newpath", map[string]string{"A": "1"})
	if err != nil {
		t.Fatalf("Changed: %v", err)
	}
	if !changed {
		t.Error("expected changed=true for uncached path")
	}
}

func TestChanged_SameSecrets(t *testing.T) {
	s := tempStore(t)
	secrets := map[string]string{"KEY": "value"}
	_ = s.Set("secret/app", secrets)
	changed, err := s.Changed("secret/app", secrets)
	if err != nil {
		t.Fatalf("Changed: %v", err)
	}
	if changed {
		t.Error("expected changed=false for identical secrets")
	}
}

func TestChanged_DifferentSecrets(t *testing.T) {
	s := tempStore(t)
	_ = s.Set("secret/app", map[string]string{"KEY": "old"})
	changed, err := s.Changed("secret/app", map[string]string{"KEY": "new"})
	if err != nil {
		t.Fatalf("Changed: %v", err)
	}
	if !changed {
		t.Error("expected changed=true for different secrets")
	}
}
