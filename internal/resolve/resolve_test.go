package resolve_test

import (
	"os"
	"testing"

	"github.com/your-org/vaultpull/internal/resolve"
)

func TestResolve_NoVariables(t *testing.T) {
	r := resolve.New(nil)
	got, err := r.Resolve("secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "secret/data/myapp" {
		t.Errorf("expected unchanged path, got %q", got)
	}
}

func TestResolve_SubstitutesFromEnvMap(t *testing.T) {
	r := resolve.New(map[string]string{"ENV": "production", "APP": "myapp"})
	got, err := r.Resolve("secret/data/${ENV}/${APP}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "secret/data/production/myapp" {
		t.Errorf("got %q", got)
	}
}

func TestResolve_FallsBackToOsEnv(t *testing.T) {
	t.Setenv("VAULT_NS", "staging")
	r := resolve.New(nil)
	got, err := r.Resolve("secret/${VAULT_NS}/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "secret/staging/db" {
		t.Errorf("got %q", got)
	}
}

func TestResolve_MissingVariable_ReturnsError(t *testing.T) {
	os.Unsetenv("MISSING_VAR")
	r := resolve.New(nil)
	_, err := r.Resolve("secret/${MISSING_VAR}/key")
	if err == nil {
		t.Fatal("expected error for missing variable")
	}
}

func TestResolveAll_AllResolved(t *testing.T) {
	r := resolve.New(map[string]string{"TIER": "prod"})
	paths := []string{"secret/${TIER}/a", "secret/${TIER}/b"}
	got, err := r.ResolveAll(paths)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0] != "secret/prod/a" || got[1] != "secret/prod/b" {
		t.Errorf("unexpected results: %v", got)
	}
}

func TestResolveAll_PropagatesError(t *testing.T) {
	os.Unsetenv("NOPE")
	r := resolve.New(nil)
	_, err := r.ResolveAll([]string{"secret/${NOPE}/x"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestHasVariables(t *testing.T) {
	if resolve.HasVariables("secret/plain/path") {
		t.Error("expected false for path without variables")
	}
	if !resolve.HasVariables("secret/${ENV}/path") {
		t.Error("expected true for path with variables")
	}
}
