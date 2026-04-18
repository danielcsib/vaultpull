package truncate_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpull/internal/truncate"
)

func TestNew_DefaultsOnZero(t *testing.T) {
	tr := truncate.New(0, "")
	if tr == nil {
		t.Fatal("expected non-nil Truncator")
	}
}

func TestValue_WithinLimit(t *testing.T) {
	tr := truncate.New(10, "...")
	got := tr.Value("hello")
	if got != "hello" {
		t.Errorf("expected 'hello', got %q", got)
	}
}

func TestValue_ExceedsLimit(t *testing.T) {
	tr := truncate.New(10, "...")
	input := "hello world!"
	got := tr.Value(input)
	if len(got) > 10 {
		t.Errorf("expected len <= 10, got %d", len(got))
	}
	if !strings.HasSuffix(got, "...") {
		t.Errorf("expected suffix '...', got %q", got)
	}
}

func TestValue_SuffixLongerThanMax(t *testing.T) {
	tr := truncate.New(2, "...")
	got := tr.Value("abcdef")
	// cutoff becomes 0, result is just suffix trimmed to nothing + suffix
	if got == "" {
		t.Error("expected non-empty result")
	}
}

func TestApply_NoTruncation(t *testing.T) {
	tr := truncate.New(100, "...")
	secrets := map[string]string{"KEY": "short"}
	out, truncated := tr.Apply(secrets)
	if out["KEY"] != "short" {
		t.Errorf("unexpected value %q", out["KEY"])
	}
	if len(truncated) != 0 {
		t.Errorf("expected no truncations, got %v", truncated)
	}
}

func TestApply_TruncatesLongValues(t *testing.T) {
	tr := truncate.New(5, "...")
	secrets := map[string]string{
		"SHORT": "hi",
		"LONG":  "this is a very long secret value",
	}
	out, truncated := tr.Apply(secrets)
	if out["SHORT"] != "hi" {
		t.Errorf("SHORT should be unchanged, got %q", out["SHORT"])
	}
	if len(out["LONG"]) > 5 {
		t.Errorf("LONG should be truncated, got len %d", len(out["LONG"]))
	}
	if len(truncated) != 1 {
		t.Errorf("expected 1 truncated key, got %d", len(truncated))
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	tr := truncate.New(3, "...")
	orig := map[string]string{"K": "longvalue"}
	origVal := orig["K"]
	tr.Apply(orig)
	if orig["K"] != origVal {
		t.Error("Apply must not mutate the input map")
	}
}
