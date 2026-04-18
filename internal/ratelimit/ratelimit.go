package ratelimit

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Limiter enforces a maximum number of operations per second.
type Limiter struct {
	mu       sync.Mutex
	rate     int
	tokens   int
	last     time.Time
	interval time.Duration
}

// New creates a Limiter that allows up to ratePerSec operations per second.
func New(ratePerSec int) (*Limiter, error) {
	if ratePerSec <= 0 {
		return nil, fmt.Errorf("ratelimit: rate must be greater than zero")
	}
	return &Limiter{
		rate:     ratePerSec,
		tokens:   ratePerSec,
		last:     time.Now(),
		interval: time.Second,
	}, nil
}

// Wait blocks until a token is available or ctx is cancelled.
func (l *Limiter) Wait(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("ratelimit: context cancelled: %w", err)
		}
		l.mu.Lock()
		l.refill()
		if l.tokens > 0 {
			l.tokens--
			l.mu.Unlock()
			return nil
		}
		l.mu.Unlock()
		select {
		case <-ctx.Done():
			return fmt.Errorf("ratelimit: context cancelled: %w", ctx.Err())
		case <-time.After(10 * time.Millisecond):
		}
	}
}

// refill adds tokens based on elapsed time. Must be called with l.mu held.
func (l *Limiter) refill() {
	now := time.Now()
	elapsed := now.Sub(l.last)
	added := int(elapsed.Seconds() * float64(l.rate))
	if added > 0 {
		l.tokens += added
		if l.tokens > l.rate {
			l.tokens = l.rate
		}
		l.last = now
	}
}
