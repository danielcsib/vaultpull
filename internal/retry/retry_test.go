package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/retry"
)

var errTemp = errors.New("temporary failure")

func TestDo_SuccessFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), retry.DefaultPolicy(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesOnFailure(t *testing.T) {
	calls := 0
	p := retry.Policy{MaxAttempts: 3, Delay: time.Millisecond, Multiplier: 1}
	err := retry.Do(context.Background(), p, func() error {
		calls++
		if calls < 3 {
			return errTemp
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	p := retry.Policy{MaxAttempts: 2, Delay: time.Millisecond, Multiplier: 1}
	err := retry.Do(context.Background(), p, func() error {
		return errTemp
	})
	if !errors.Is(err, retry.ErrExhausted) {
		t.Fatalf("expected ErrExhausted, got %v", err)
	}
	if !errors.Is(err, errTemp) {
		t.Fatalf("expected wrapped errTemp, got %v", err)
	}
}

func TestDo_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	p := retry.Policy{MaxAttempts: 5, Delay: time.Millisecond, Multiplier: 1}
	err := retry.Do(ctx, p, func() error { return errTemp })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestDo_ZeroMaxAttempts(t *testing.T) {
	calls := 0
	p := retry.Policy{MaxAttempts: 0, Delay: time.Millisecond}
	_ = retry.Do(context.Background(), p, func() error {
		calls++
		return nil
	})
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}
