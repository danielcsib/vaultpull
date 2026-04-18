package retry

import (
	"context"
	"errors"
	"time"
)

// Policy defines retry behaviour.
type Policy struct {
	MaxAttempts int
	Delay       time.Duration
	Multiplier  float64
}

// DefaultPolicy returns a sensible default retry policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts: 3,
		Delay:       500 * time.Millisecond,
		Multiplier:  2.0,
	}
}

// ErrExhausted is returned when all attempts are consumed.
var ErrExhausted = errors.New("retry: all attempts exhausted")

// Do runs fn up to MaxAttempts times, backing off between failures.
// It respects context cancellation.
func Do(ctx context.Context, p Policy, fn func() error) error {
	if p.MaxAttempts <= 0 {
		p.MaxAttempts = 1
	}
	delay := p.Delay
	var last error
	for i := 0; i < p.MaxAttempts; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		last = fn()
		if last == nil {
			return nil
		}
		if i < p.MaxAttempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			if p.Multiplier > 1 {
				delay = time.Duration(float64(delay) * p.Multiplier)
			}
		}
	}
	return errors.Join(ErrExhausted, last)
}
