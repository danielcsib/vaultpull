package envlock_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourusername/vaultpull/internal/envlock"
)

func ExampleLocker_Acquire() {
	dir, _ := os.MkdirTemp("", "envlock-example")
	defer os.RemoveAll(dir)

	envPath := filepath.Join(dir, ".env")
	l := envlock.New(envPath, 0)

	if err := l.Acquire(); err != nil {
		if errors.Is(err, envlock.ErrLockHeld) {
			fmt.Println("lock held by another process")
			return
		}
		panic(err)
	}
	defer l.Release()

	fmt.Println("lock acquired")
	fmt.Println("lock held:", l.Held())
	// Output:
	// lock acquired
	// lock held: true
}
