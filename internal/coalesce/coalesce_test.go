package coalesce_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/coalesce"
)

func TestMerge_NoSources(t *testing.T) {
	m := coalesce.New()
	_, err := m.Merge()
	if err == nil {
		t.Fatal("expected error for no sources, got nil")
	}
}

func TestMerge_SingleSource(t *testing.T) {
	m := coalesce.New()
	src := map[string]string{"KEY": "value"}
	out, err := m.Merge(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "value" {
		t.Errorf("expected 'value', got %q", out["KEY"])
	}
}

func TestMerge_HigherPriorityWins(t *testing.T) {
	m := coalesce.New()
	high := map[string]string{"DB_PASS": "secret1"}
	low := map[string]string{"DB_PASS": "secret2", "API_KEY": "abc"}
	out, err := m.Merge(high, low)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_PASS"] != "secret1" {
		t.Errorf("expected 'secret1', got %q", out["DB_PASS"])
	}
	if out["API_KEY"] != "abc" {
		t.Errorf("expected 'abc', got %q", out["API_KEY"])
	}
}

func TestMerge_EmptyValueWinsWithoutOmitEmpty(t *testing.T) {
	m := coalesce.New()
	high := map[string]string{"KEY": ""}
	low := map[string]string{"KEY": "fallback"}
	out, err := m.Merge(high, low)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "" {
		t.Errorf("expected empty string, got %q", out["KEY"])
	}
}

func TestMerge_OmitEmptyFallsThrough(t *testing.T) {
	m := coalesce.New(coalesce.OmitEmpty())
	high := map[string]string{"KEY": ""}
	low := map[string]string{"KEY": "fallback"}
	out, err := m.Merge(high, low)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "fallback" {
		t.Errorf("expected 'fallback', got %q", out["KEY"])
	}
}

func TestMerge_OmitEmpty_AllEmpty_KeyAbsent(t *testing.T) {
	m := coalesce.New(coalesce.OmitEmpty())
	high := map[string]string{"KEY": ""}
	low := map[string]string{"KEY": ""}
	out, err := m.Merge(high, low)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["KEY"]; ok {
		t.Error("expected KEY to be absent from output")
	}
}

func TestFirst_ReturnsFirstMatch(t *testing.T) {
	a := map[string]string{"X": "1"}
	b := map[string]string{"X": "2", "Y": "3"}
	v, ok := coalesce.First("X", a, b)
	if !ok || v != "1" {
		t.Errorf("expected '1', got %q ok=%v", v, ok)
	}
}

func TestFirst_FallsToSecondSource(t *testing.T) {
	a := map[string]string{}
	b := map[string]string{"Y": "found"}
	v, ok := coalesce.First("Y", a, b)
	if !ok || v != "found" {
		t.Errorf("expected 'found', got %q ok=%v", v, ok)
	}
}

func TestFirst_NotFound(t *testing.T) {
	_, ok := coalesce.First("MISSING", map[string]string{})
	if ok {
		t.Error("expected not found")
	}
}
