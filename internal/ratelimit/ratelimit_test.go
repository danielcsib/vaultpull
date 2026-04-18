package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/ratelimit"
)

func TestNew_InvalidRate(t *testing.T) {
	_, err := ratelimit.New(0)
	if err == nil {
		t.Fatal("expected error for zero rate")
	}
}

func TestNew_ValidRate(t *testing.T) {
	l, err := ratelimit.New(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
}

func TestWait_AllowsUpToRate(t *testing.T) {
	l, _ := ratelimit.New(5)
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("unexpected error on call %d: %v", i, err)
		}
	}
}

func TestWait_ContextCancelled(t *testing.T) {
	l, _ := ratelimit.New(1)
	ctx := context.Background()
	// consume the single token
	_ = l.Wait(ctx)

	ctx2, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	err := l.Wait(ctx2)
	if err == nil {
		t.Fatal("expected error when context times out")
	}
}

func TestWait_RefillsOverTime(t *testing.T) {
	l, _ := ratelimit.New(50)
	ctx := context.Background()
	// drain all tokens
	for i := 0; i < 50; i++ {
		_ = l.Wait(ctx)
	}
	// wait for refill
	time.Sleep(100 * time.Millisecond)
	if err := l.Wait(ctx); err != nil {
		t.Fatalf("expected token after refill: %v", err)
	}
}
