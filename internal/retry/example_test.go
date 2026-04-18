package retry_test

import (
	"context"
	"fmt"
	"time"

	"github.com/your-org/vaultpull/internal/retry"
)

func ExampleDo() {
	attempts := 0
	p := retry.Policy{
		MaxAttempts: 3,
		Delay:       time.Millisecond,
		Multiplier:  1.0,
	}
	err := retry.Do(context.Background(), p, func() error {
		attempts++
		if attempts < 2 {
			return fmt.Errorf("not ready")
		}
		return nil
	})
	fmt.Println(err, attempts)
	// Output: <nil> 2
}
