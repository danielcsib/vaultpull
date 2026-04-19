package dedupe_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/dedupe"
)

func TestMerge_NoConflict(t *testing.T) {
	m := dedupe.New(dedupe.PolicyKeepFirst)
	dst := map[string]string{"A": "1"}
	src := map[string]string{"B": "2"}
	if err := m.Merge(dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["B"] != "2" {
		t.Errorf("expected B=2, got %q", dst["B"])
	}
	if len(m.Conflicts) != 0 {
		t.Errorf("expected no conflicts")
	}
}

func TestMerge_PolicyKeepFirst(t *testing.T) {
	m := dedupe.New(dedupe.PolicyKeepFirst)
	dst := map[string]string{"KEY": "original"}
	src := map[string]string{"KEY": "new"}
	_ = m.Merge(dst, src)
	if dst["KEY"] != "original" {
		t.Errorf("expected original, got %q", dst["KEY"])
	}
	if len(m.Conflicts) != 1 {
		t.Errorf("expected 1 conflict")
	}
}

func TestMerge_PolicyKeepLast(t *testing.T) {
	m := dedupe.New(dedupe.PolicyKeepLast)
	dst := map[string]string{"KEY": "original"}
	src := map[string]string{"KEY": "new"}
	_ = m.Merge(dst, src)
	if dst["KEY"] != "new" {
		t.Errorf("expected new, got %q", dst["KEY"])
	}
}

func TestMerge_PolicyError(t *testing.T) {
	m := dedupe.New(dedupe.PolicyError)
	dst := map[string]string{"KEY": "original"}
	src := map[string]string{"KEY": "new"}
	err := m.Merge(dst, src)
	if err == nil {
		t.Fatal("expected error on conflict")
	}
}

func TestReset_ClearsConflicts(t *testing.T) {
	m := dedupe.New(dedupe.PolicyKeepFirst)
	dst := map[string]string{"X": "a"}
	src := map[string]string{"X": "b"}
	_ = m.Merge(dst, src)
	m.Reset()
	if len(m.Conflicts) != 0 {
		t.Errorf("expected conflicts cleared")
	}
}
