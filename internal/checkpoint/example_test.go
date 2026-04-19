package checkpoint_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourusername/vaultpull/internal/checkpoint"
)

func ExampleStore_Record() {
	dir, _ := os.MkdirTemp("", "cp")
	defer os.RemoveAll(dir)

	store, _ := checkpoint.NewStore(filepath.Join(dir, "checkpoint.json"))
	_ = store.Record("secret/myapp", 4)

	entry, ok := store.Get("secret/myapp")
	if ok {
		fmt.Println("key count:", entry.KeyCount)
	}
	// Output:
	// key count: 4
}
