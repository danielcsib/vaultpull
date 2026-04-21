package quota_test

import (
	"errors"
	"testing"

	"github.com/your-org/vaultpull/internal/quota"
)

func TestNew_DefaultMaxReads(t *testing.T) {
	s := quota.New(0)
	if s == nil {
		t.Fatal("expected non-nil store")
	}
	// default is 10; 10 reads should succeed
	for i := 0; i < 10; i++ {
		if err := s.Check("secret/foo"); err != nil {
			t.Fatalf("unexpected error on read %d: %v", i+1, err)
		}
	}
}

func TestCheck_WithinQuota(t *testing.T) {
	s := quota.New(3)
	for i := 0; i < 3; i++ {
		if err := s.Check("secret/bar"); err != nil {
			t.Fatalf("read %d should be allowed: %v", i+1, err)
		}
	}
}

func TestCheck_ExceedsQuota(t *testing.T) {
	s := quota.New(2)
	_ = s.Check("secret/x")
	_ = s.Check("secret/x")
	err := s.Check("secret/x")
	if !errors.Is(err, quota.ErrQuotaExceeded) {
		t.Fatalf("expected ErrQuotaExceeded, got %v", err)
	}
}

func TestRemaining_DecrementsCorrectly(t *testing.T) {
	s := quota.New(5)
	if got := s.Remaining("secret/y"); got != 5 {
		t.Fatalf("expected 5 remaining, got %d", got)
	}
	_ = s.Check("secret/y")
	_ = s.Check("secret/y")
	if got := s.Remaining("secret/y"); got != 3 {
		t.Fatalf("expected 3 remaining, got %d", got)
	}
}

func TestRemaining_NeverNegative(t *testing.T) {
	s := quota.New(1)
	_ = s.Check("secret/z")
	_ = s.Check("secret/z") // exceeds
	if got := s.Remaining("secret/z"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestReset_ClearsPath(t *testing.T) {
	s := quota.New(1)
	_ = s.Check("secret/a")
	s.Reset("secret/a")
	if err := s.Check("secret/a"); err != nil {
		t.Fatalf("expected read to succeed after reset: %v", err)
	}
}

func TestResetAll_ClearsAll(t *testing.T) {
	s := quota.New(1)
	_ = s.Check("secret/a")
	_ = s.Check("secret/b")
	s.ResetAll()
	snap := s.Snapshot()
	if len(snap) != 0 {
		t.Fatalf("expected empty snapshot after ResetAll, got %v", snap)
	}
}

func TestSnapshot_IsCopy(t *testing.T) {
	s := quota.New(5)
	_ = s.Check("secret/snap")
	snap := s.Snapshot()
	snap["secret/snap"] = 999
	if got := s.Remaining("secret/snap"); got != 4 {
		t.Fatalf("snapshot mutation should not affect store, remaining=%d", got)
	}
}
