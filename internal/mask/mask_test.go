package mask_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/mask"
)

func TestIsSensitive_MatchesDefault(t *testing.T) {
	m := mask.New()
	cases := []struct {
		key      string
		want     bool
	}{
		{"DB_PASSWORD", true},
		{"API_TOKEN", true},
		{"SECRET_KEY", true},
		{"AUTH_HEADER", true},
		{"DATABASE_URL", false},
		{"APP_ENV", false},
	}
	for _, tc := range cases {
		got := m.IsSensitive(tc.key)
		if got != tc.want {
			t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.want)
		}
	}
}

func TestMaskMap_RedactsSensitive(t *testing.T) {
	m := mask.New()
	input := map[string]string{
		"DB_PASSWORD": "supersecret",
		"APP_ENV":     "production",
		"API_TOKEN":   "tok_abc123",
	}
	out := m.MaskMap(input)
	if out["DB_PASSWORD"] != "***REDACTED***" {
		t.Errorf("expected DB_PASSWORD to be redacted, got %q", out["DB_PASSWORD"])
	}
	if out["API_TOKEN"] != "***REDACTED***" {
		t.Errorf("expected API_TOKEN to be redacted, got %q", out["API_TOKEN"])
	}
	if out["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV to be unchanged, got %q", out["APP_ENV"])
	}
}

func TestMaskMap_DoesNotMutateInput(t *testing.T) {
	m := mask.New()
	input := map[string]string{"DB_PASSWORD": "supersecret"}
	_ = m.MaskMap(input)
	if input["DB_PASSWORD"] != "supersecret" {
		t.Error("MaskMap mutated the input map")
	}
}

func TestNewWithPatterns_CustomPlaceholder(t *testing.T) {
	m := mask.NewWithPatterns([]string{"private"}, "[hidden]")
	out := m.MaskMap(map[string]string{
		"PRIVATE_KEY": "abc",
		"PUBLIC_URL":  "https://example.com",
	})
	if out["PRIVATE_KEY"] != "[hidden]" {
		t.Errorf("expected [hidden], got %q", out["PRIVATE_KEY"])
	}
	if out["PUBLIC_URL"] != "https://example.com" {
		t.Errorf("expected PUBLIC_URL unchanged, got %q", out["PUBLIC_URL"])
	}
}
