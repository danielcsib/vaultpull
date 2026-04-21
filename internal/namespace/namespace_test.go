package namespace_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/namespace"
)

func TestNew_EmptyNamespace(t *testing.T) {
	_, err := namespace.New("", "/")
	if err == nil {
		t.Fatal("expected error for empty namespace")
	}
}

func TestNew_DefaultSeparator(t *testing.T) {
	s, err := namespace.New("prod", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := s.Qualify("DB_PASS")
	want := "prod/DB_PASS"
	if got != want {
		t.Errorf("Qualify = %q, want %q", got, want)
	}
}

func TestQualify_AddsPrefix(t *testing.T) {
	s, _ := namespace.New("staging", ":")
	got := s.Qualify("API_KEY")
	if got != "staging:API_KEY" {
		t.Errorf("got %q", got)
	}
}

func TestStrip_RemovesPrefix(t *testing.T) {
	s, _ := namespace.New("staging", ":")
	got := s.Strip("staging:API_KEY")
	if got != "API_KEY" {
		t.Errorf("got %q", got)
	}
}

func TestStrip_NoPrefix_Passthrough(t *testing.T) {
	s, _ := namespace.New("prod", "/")
	got := s.Strip("OTHER_KEY")
	if got != "OTHER_KEY" {
		t.Errorf("got %q", got)
	}
}

func TestApply_QualifiesAllKeys(t *testing.T) {
	s, _ := namespace.New("prod", "/")
	input := map[string]string{"FOO": "1", "BAR": "2"}
	out := s.Apply(input)
	for _, k := range []string{"prod/FOO", "prod/BAR"} {
		if _, ok := out[k]; !ok {
			t.Errorf("missing key %q in output", k)
		}
	}
	if len(out) != len(input) {
		t.Errorf("length mismatch: got %d, want %d", len(out), len(input))
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	s, _ := namespace.New("prod", "/")
	input := map[string]string{"FOO": "bar"}
	s.Apply(input)
	if _, ok := input["prod/FOO"]; ok {
		t.Error("Apply mutated the input map")
	}
}

func TestUnwrap_StripsAllKeys(t *testing.T) {
	s, _ := namespace.New("prod", "/")
	input := map[string]string{"prod/FOO": "1", "prod/BAR": "2"}
	out := s.Unwrap(input)
	for _, k := range []string{"FOO", "BAR"} {
		if _, ok := out[k]; !ok {
			t.Errorf("missing key %q after unwrap", k)
		}
	}
}

func TestPrefix_ReturnsNamespace(t *testing.T) {
	s, _ := namespace.New("dev", "/")
	if s.Prefix() != "dev" {
		t.Errorf("Prefix = %q, want %q", s.Prefix(), "dev")
	}
}
