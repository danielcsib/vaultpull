package selector_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/selector"
)

func TestNew_EmptyKeys(t *testing.T) {
	_, err := selector.New([]string{})
	if err == nil {
		t.Fatal("expected error for empty keys, got nil")
	}
}

func TestNew_EmptyStringKey(t *testing.T) {
	_, err := selector.New([]string{"VALID", ""})
	if err == nil {
		t.Fatal("expected error for empty string key, got nil")
	}
}

func TestNew_ValidKeys(t *testing.T) {
	s, err := selector.New([]string{"DB_HOST", "DB_PORT"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := s.Keys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
}

func TestNew_DeduplicatesKeys(t *testing.T) {
	s, err := selector.New([]string{"KEY", "KEY", "KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Keys()) != 1 {
		t.Fatalf("expected 1 unique key, got %d", len(s.Keys()))
	}
}

func TestPick_ReturnsOnlySelectedKeys(t *testing.T) {
	s, _ := selector.New([]string{"A", "B"})
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}

	picked, missing := s.Pick(secrets)

	if len(picked) != 2 {
		t.Fatalf("expected 2 picked keys, got %d", len(picked))
	}
	if picked["A"] != "1" || picked["B"] != "2" {
		t.Errorf("unexpected picked values: %v", picked)
	}
	if _, ok := picked["C"]; ok {
		t.Error("key C should not be in picked map")
	}
	if len(missing) != 0 {
		t.Errorf("expected no missing keys, got %v", missing)
	}
}

func TestPick_ReportsMissingKeys(t *testing.T) {
	s, _ := selector.New([]string{"PRESENT", "ABSENT"})
	secrets := map[string]string{"PRESENT": "yes"}

	_, missing := s.Pick(secrets)

	if len(missing) != 1 || missing[0] != "ABSENT" {
		t.Errorf("expected missing=[ABSENT], got %v", missing)
	}
}

func TestHas_ReturnsTrueForKnownKey(t *testing.T) {
	s, _ := selector.New([]string{"TOKEN"})
	if !s.Has("TOKEN") {
		t.Error("expected Has to return true for TOKEN")
	}
	if s.Has("OTHER") {
		t.Error("expected Has to return false for OTHER")
	}
}

func TestKeys_ReturnsSorted(t *testing.T) {
	s, _ := selector.New([]string{"Z", "A", "M"})
	keys := s.Keys()
	expected := []string{"A", "M", "Z"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("position %d: expected %s, got %s", i, expected[i], k)
		}
	}
}
