package timeout

import (
	"context"
	"fmt"
	"time"
)

// DefaultTimeout is used when no timeout is specified.
const DefaultTimeout = 10 * time.Second

// Config holds timeout configuration.
type Config struct {
	Duration time.Duration
}

// New returns a Config with the given duration.
// If d is zero or negative, DefaultTimeout is used.
func New(d time.Duration) Config {
	if d <= 0 {
		d = DefaultTimeout
	}
	return Config{Duration: d}
}

// WithContext returns a context that cancels after the configured duration.
// The caller is responsible for calling the returned cancel function.
func (c Config) WithContext(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, c.Duration)
}

// Wrap executes fn within a timeout context derived from parent.
// Returns an error if fn does not complete within the configured duration.
func (c Config) Wrap(parent context.Context, fn func(ctx context.Context) error) error {
	ctx, cancel := c.WithContext(parent)
	defer cancel()

	type result struct {
		err error
	}
	ch := make(chan result, 1)
	go func() {
		ch <- result{err: fn(ctx)}
	}()

	select {
	case r := <-ch:
		return r.err
	case <-ctx.Done():
		return fmt.Errorf("operation timed out after %s", c.Duration)
	}
}
