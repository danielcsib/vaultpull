package watch_test

import (
	"context"
	"fmt"
	"time"

	"github.com/your-org/vaultpull/internal/watch"
)

type staticClient struct{}

func (s *staticClient) ReadSecrets(_ string) (map[string]string, error) {
	return map[string]string{"DB_PASS": "hunter2"}, nil
}

func ExampleNew() {
	client := &staticClient{}

	var once bool
	w := watch.New(client, []string{"secret/app"}, 50*time.Millisecond,
		func(path string, secrets map[string]string) {
			if !once {
				fmt.Printf("changed: %s keys=%d\n", path, len(secrets))
				once = true
			}
		},
	)

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	w.Run(ctx)

	// Output:
	// changed: secret/app keys=1
}
