package sanitize_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/sanitize"
)

func TestValue_TrimSpace(t *testing.T) {
	s := sanitize.New(sanitize.Options{TrimSpace: true})
	got := s.Value("  hello  ")
	if got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}

func TestValue_NoTrimSpace(t *testing.T) {
	s := sanitize.New(sanitize.Options{TrimSpace: false})
	got := s.Value("  hello  ")
	if got != "  hello  " {
		t.Fatalf("expected %q, got %q", "  hello  ", got)
	}
}

func TestValue_StripNonPrintable(t *testing.T) {
	s := sanitize.New(sanitize.Options{StripNonPrintable: true})
	input := "hello\x00world\x01"
	got := s.Value(input)
	if got != "helloworld" {
		t.Fatalf("expected %q, got %q", "helloworld", got)
	}
}

func TestValue_PreservesTabAndNewline(t *testing.T) {
	s := sanitize.New(sanitize.Options{StripNonPrintable: true})
	input := "line1\nline2\ttabbed"
	got := s.Value(input)
	if got != input {
		t.Fatalf("expected %q, got %q", input, got)
	}
}

func TestValue_NormaliseNewlines(t *testing.T) {
	s := sanitize.New(sanitize.Options{NormaliseNewlines: true})
	got := s.Value("a\r\nb\rc")
	if got != "a\nb\nc" {
		t.Fatalf("expected %q, got %q", "a\nb\nc", got)
	}
}

func TestApply_ReturnsNewMap(t *testing.T) {
	s := sanitize.New(sanitize.DefaultOptions())
	input := map[string]string{
		"KEY_A": "  value  ",
		"KEY_B": "clean",
	}
	out := s.Apply(input)
	if out["KEY_A"] != "value" {
		t.Fatalf("KEY_A: expected %q, got %q", "value", out["KEY_A"])
	}
	if out["KEY_B"] != "clean" {
		t.Fatalf("KEY_B: expected %q, got %q", "clean", out["KEY_B"])
	}
	// Original must not be mutated.
	if input["KEY_A"] != "  value  " {
		t.Fatal("Apply mutated the original map")
	}
}

func TestApply_EmptyMap(t *testing.T) {
	s := sanitize.New(sanitize.DefaultOptions())
	out := s.Apply(map[string]string{})
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %v", out)
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := sanitize.DefaultOptions()
	if !opts.TrimSpace || !opts.StripNonPrintable || !opts.NormaliseNewlines {
		t.Fatal("DefaultOptions should enable all sanitisation steps")
	}
}
