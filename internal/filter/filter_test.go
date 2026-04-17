package filter_test

import (
	"testing"

	"github.com/example/vaultpull/internal/filter"
)

func TestNew_EmptyPatterns(t *testing.T) {
	f := filter.New([]string{})
	if !f.Allow("ANY_KEY") {
		t.Error("expected all keys allowed when no rules defined")
	}
}

func TestAllow_IncludePrefix(t *testing.T) {
	f := filter.New([]string{"APP_"})
	if !f.Allow("APP_SECRET") {
		t.Error("expected APP_SECRET to be allowed")
	}
	if f.Allow("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD to be excluded")
	}
}

func TestAllow_ExcludePrefix(t *testing.T) {
	f := filter.New([]string{"!INTERNAL_"})
	if !f.Allow("APP_KEY") {
		t.Error("expected APP_KEY to be allowed")
	}
	if f.Allow("INTERNAL_SECRET") {
		t.Error("expected INTERNAL_SECRET to be excluded")
	}
}

func TestAllow_ExcludeTakesPrecedence(t *testing.T) {
	f := filter.New([]string{"APP_", "!APP_INTERNAL_"})
	if !f.Allow("APP_PUBLIC") {
		t.Error("expected APP_PUBLIC allowed")
	}
	if f.Allow("APP_INTERNAL_KEY") {
		t.Error("expected APP_INTERNAL_KEY excluded")
	}
}

func TestApply_FiltersMap(t *testing.T) {
	f := filter.New([]string{"DB_"})
	secrets := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "secret",
		"APP_TOKEN":   "abc123",
	}
	result := f.Apply(secrets)
	if len(result) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["APP_TOKEN"]; ok {
		t.Error("APP_TOKEN should have been filtered out")
	}
}

func TestApply_EmptyFilter_AllowsAll(t *testing.T) {
	f := filter.New(nil)
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	result := f.Apply(secrets)
	if len(result) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(result))
	}
}
