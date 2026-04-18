package validate_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/validate"
)

func TestCheck_AllPresent(t *testing.T) {
	secrets := map[string]string{"DB_URL": "postgres://localhost", "API_KEY": "abc123"}
	r := validate.Check(secrets, []string{"DB_URL", "API_KEY"})
	if r.HasErrors() {
		t.Fatalf("expected no errors, got missing: %v", r.Missing)
	}
}

func TestCheck_MissingKey(t *testing.T) {
	secrets := map[string]string{"DB_URL": "postgres://localhost"}
	r := validate.Check(secrets, []string{"DB_URL", "API_KEY"})
	if !r.HasErrors() {
		t.Fatal("expected errors for missing key")
	}
	if len(r.Missing) != 1 || r.Missing[0] != "API_KEY" {
		t.Fatalf("unexpected missing keys: %v", r.Missing)
	}
}

func TestCheck_EmptyValue(t *testing.T) {
	secrets := map[string]string{"DB_URL": ""}
	r := validate.Check(secrets, []string{"DB_URL"})
	if r.HasErrors() {
		t.Fatal("empty value should not be a hard error")
	}
	if len(r.Empty) != 1 || r.Empty[0] != "DB_URL" {
		t.Fatalf("expected empty key reported, got: %v", r.Empty)
	}
}

func TestCheck_WhitespaceWarning(t *testing.T) {
	secrets := map[string]string{"TOKEN": "  secret  "}
	r := validate.Check(secrets, nil)
	if len(r.Warnings) == 0 {
		t.Fatal("expected whitespace warning")
	}
}

func TestCheck_NoRequiredKeys(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	r := validate.Check(secrets, nil)
	if r.HasErrors() {
		t.Fatal("expected no errors with no required keys")
	}
}

func TestResult_Summary(t *testing.T) {
	r := validate.Result{
		Missing:  []string{"SECRET_KEY"},
		Empty:    []string{"OPTIONAL"},
		Warnings: []string{"key \"X\" has leading/trailing whitespace"},
	}
	s := r.Summary()
	if !strings.Contains(s, "[ERROR]") {
		t.Error("summary should contain ERROR")
	}
	if !strings.Contains(s, "[WARN]") {
		t.Error("summary should contain WARN")
	}
}
