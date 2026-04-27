// Package envlock provides file-based locking for .env files to prevent
// concurrent writes from corrupting secrets during sync operations.
package envlock

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ErrLockHeld is returned when a lock file already exists and has not expired.
var ErrLockHeld = errors.New("envlock: lock is already held")

// DefaultTTL is the maximum age of a lock file before it is considered stale.
const DefaultTTL = 30 * time.Second

// Locker manages an advisory lock file for a given .env path.
type Locker struct {
	path string
	ttl  time.Duration
}

// New creates a Locker for the given env file path.
// The lock file is placed alongside the target as "<file>.lock".
func New(envPath string, ttl time.Duration) *Locker {
	if ttl <= 0 {
		ttl = DefaultTTL
	}
	return &Locker{
		path: lockPath(envPath),
		ttl:  ttl,
	}
}

// Acquire creates the lock file. Returns ErrLockHeld if a valid lock exists.
func (l *Locker) Acquire() error {
	info, err := os.Stat(l.path)
	if err == nil {
		if time.Since(info.ModTime()) < l.ttl {
			return fmt.Errorf("%w: %s", ErrLockHeld, l.path)
		}
		// Stale lock — remove it before acquiring.
		_ = os.Remove(l.path)
	}

	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("%w: %s", ErrLockHeld, l.path)
		}
		return fmt.Errorf("envlock: create lock: %w", err)
	}
	f.Close()
	return nil
}

// Release removes the lock file. It is safe to call even if the lock is not held.
func (l *Locker) Release() error {
	err := os.Remove(l.path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("envlock: release lock: %w", err)
	}
	return nil
}

// Held reports whether the lock file exists and has not expired.
func (l *Locker) Held() bool {
	info, err := os.Stat(l.path)
	if err != nil {
		return false
	}
	return time.Since(info.ModTime()) < l.ttl
}

func lockPath(envPath string) string {
	return filepath.Join(filepath.Dir(envPath), filepath.Base(envPath)+".lock")
}
