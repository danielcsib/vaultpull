package rollback_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/your-org/vaultpull/internal/rollback"
)

func ExampleRollbacker_Restore() {
	dir, _ := os.MkdirTemp("", "rb-example-*")
	defer os.RemoveAll(dir)

	// Write a fake backup.
	backup := filepath.Join(dir, ".env.20240101T000000.bak")
	os.WriteFile(backup, []byte("API_KEY=secret\n"), 0600)

	target := filepath.Join(dir, ".env")
	os.WriteFile(target, []byte("API_KEY=old\n"), 0600)

	rb := rollback.New(dir)
	src, err := rb.Restore(target)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	data, _ := os.ReadFile(target)
	fmt.Printf("restored from %s\n", filepath.Base(src))
	fmt.Printf("content: %s", data)
	// Output:
	// restored from .env.20240101T000000.bak
	// content: API_KEY=secret
}
