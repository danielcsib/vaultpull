package redact

import (
	"strings"
	"testing"
)

func TestMask_EmptyValue(t *testing.T) {
	r := New(4)
	if got := r.Mask(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestMask_ZeroVisible(t *testing.T) {
	r := New(0)
	got := r.Mask("secret")
	if got != "******" {
		t.Fatalf("expected all stars, got %q", got)
	}
}

func TestMask_PartialReveal(t *testing.T) {
	r := New(4)
	got := r.Mask("mysecretvalue")
	if !strings.HasSuffix(got, "alue") {
		t.Fatalf("expected suffix 'alue', got %q", got)
	}
	if !strings.HasPrefix(got, "*") {
		t.Fatalf("expected leading stars, got %q", got)
	}
	if len(got) != len("mysecretvalue") {
		t.Fatalf("length mismatch: got %d", len(got))
	}
}

func TestMask_ShortValue(t *testing.T) {
	r := New(4)
	got := r.Mask("ab")
	if got != "**" {
		t.Fatalf("expected '**', got %q", got)
	}
}

func TestMaskMap(t *testing.T) {
	r := New(2)
	input := map[string]string{
		"KEY1": "password",
		"KEY2": "",
	}
	out := r.MaskMap(input)
	if out["KEY2"] != "" {
		t.Fatalf("expected empty for KEY2, got %q", out["KEY2"])
	}
	if !strings.HasSuffix(out["KEY1"], "rd") {
		t.Fatalf("expected suffix 'rd' for KEY1, got %q", out["KEY1"])
	}
	if _, ok := out["KEY1"]; !ok {
		t.Fatal("KEY1 missing from output")
	}
}

func TestNew_NegativeVisible(t *testing.T) {
	r := New(-5)
	if r.visibleChars != 0 {
		t.Fatalf("expected visibleChars=0, got %d", r.visibleChars)
	}
}
