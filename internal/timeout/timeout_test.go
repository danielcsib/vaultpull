package timeout_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/example/vaultpull/internal/timeout"
)

func TestNew_DefaultsOnZero(t *testing.T) {
	cfg := timeout.New(0)
	if cfg.Duration != timeout.DefaultTimeout {
		t.Fatalf("expected default timeout, got %s", cfg.Duration)
	}
}

func TestNew_CustomDuration(t *testing.T) {
	cfg := timeout.New(5 * time.Second)
	if cfg.Duration != 5*time.Second {
		t.Fatalf("expected 5s, got %s", cfg.Duration)
	}
}

func TestWrap_SuccessWithinTimeout(t *testing.T) {
	cfg := timeout.New(1 * time.Second)
	err := cfg.Wrap(context.Background(), func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWrap_PropagatesError(t *testing.T) {
	cfg := timeout.New(1 * time.Second)
	sentinel := errors.New("vault error")
	err := cfg.Wrap(context.Background(), func(ctx context.Context) error {
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestWrap_TimesOut(t *testing.T) {
	cfg := timeout.New(50 * time.Millisecond)
	err := cfg.Wrap(context.Background(), func(ctx context.Context) error {
		time.Sleep(200 * time.Millisecond)
		return nil
	})
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestWithContext_CancelReleasesResources(t *testing.T) {
	cfg := timeout.New(1 * time.Second)
	ctx, cancel := cfg.WithContext(context.Background())
	cancel()
	select {
	case <-ctx.Done():
		// expected
	default:
		t.Fatal("expected context to be cancelled")
	}
}
