package envlock_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/envlock"
)

func tempEnvFile(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, ".env")
}

func TestAcquire_CreatesLockFile(t *testing.T) {
	path := tempEnvFile(t)
	l := envlock.New(path, 0)

	if err := l.Acquire(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer l.Release()

	if _, err := os.Stat(path + ".lock"); err != nil {
		t.Errorf("lock file not created: %v", err)
	}
}

func TestAcquire_ReturnsErrLockHeld(t *testing.T) {
	path := tempEnvFile(t)
	l := envlock.New(path, 5*time.Second)

	if err := l.Acquire(); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	defer l.Release()

	err := l.Acquire()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, envlock.ErrLockHeld) {
		t.Errorf("expected ErrLockHeld, got %v", err)
	}
}

func TestRelease_RemovesLockFile(t *testing.T) {
	path := tempEnvFile(t)
	l := envlock.New(path, 0)

	_ = l.Acquire()
	if err := l.Release(); err != nil {
		t.Fatalf("release failed: %v", err)
	}

	if _, err := os.Stat(path + ".lock"); !os.IsNotExist(err) {
		t.Error("lock file still exists after release")
	}
}

func TestRelease_IdempotentWhenNotHeld(t *testing.T) {
	path := tempEnvFile(t)
	l := envlock.New(path, 0)

	if err := l.Release(); err != nil {
		t.Errorf("release on unheld lock should not error: %v", err)
	}
}

func TestHeld_ReturnsFalseWhenNoLock(t *testing.T) {
	path := tempEnvFile(t)
	l := envlock.New(path, 0)

	if l.Held() {
		t.Error("expected Held to return false")
	}
}

func TestHeld_ReturnsTrueAfterAcquire(t *testing.T) {
	path := tempEnvFile(t)
	l := envlock.New(path, 5*time.Second)

	_ = l.Acquire()
	defer l.Release()

	if !l.Held() {
		t.Error("expected Held to return true")
	}
}

func TestAcquire_StaleLockIsReplaced(t *testing.T) {
	path := tempEnvFile(t)
	l := envlock.New(path, 1*time.Millisecond)

	_ = l.Acquire()
	time.Sleep(5 * time.Millisecond)

	l2 := envlock.New(path, 5*time.Second)
	if err := l2.Acquire(); err != nil {
		t.Fatalf("expected stale lock to be replaced: %v", err)
	}
	defer l2.Release()
}
