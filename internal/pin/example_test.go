package pin_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/your-org/vaultpull/internal/pin"
)

func ExampleStore_Pin() {
	dir, _ := os.MkdirTemp("", "pin-example")
	defer os.RemoveAll(dir)

	store, _ := pin.NewStore(filepath.Join(dir, "pins.json"))

	_ = store.Pin("secret/app/db", 4, "alice")

	e, ok := store.Get("secret/app/db")
	if ok {
		fmt.Printf("pinned %s to version %d by %s\n", e.Path, e.Version, e.PinnedBy)
	}

	// Output:
	// pinned secret/app/db to version 4 by alice
}
