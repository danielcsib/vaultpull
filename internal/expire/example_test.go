package expire_test

import (
	"fmt"
	"os"
	"time"

	"github.com/your-org/vaultpull/internal/expire"
)

func ExampleStore_Set() {
	f, _ := os.CreateTemp("", "expire-*.json")
	f.Close()
	defer os.Remove(f.Name())

	store, _ := expire.NewStore(f.Name())

	_ = store.Set("secret/app", 1*time.Hour)

	if store.Expired("secret/app") {
		fmt.Println("expired")
	} else {
		fmt.Println("still valid")
	}

	// Output: still valid
}
