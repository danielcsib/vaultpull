package transform

import (
	"testing"
)

func TestApply_NilMap(t *testing.T) {
	_, err := Apply(nil, Rule{})
	if err == nil {
		t.Fatal("expected error for nil map")
	}
}

func TestApply_UppercaseKeys(t *testing.T) {
	secrets := map[string]string{"db_host": "localhost", "db_port": "5432"}
	result, err := Apply(secrets, Rule{Uppercase: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", result["DB_HOST"])
	}
	if result["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", result["DB_PORT"])
	}
}

func TestApply_PrefixKeys(t *testing.T) {
	secrets := map[string]string{"TOKEN": "abc"}
	result, err := Apply(secrets, Rule{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["APP_TOKEN"] != "abc" {
		t.Errorf("expected APP_TOKEN=abc, got %v", result)
	}
}

func TestApply_StripPrefix(t *testing.T) {
	secrets := map[string]string{"vault_secret": "val"}
	result, err := Apply(secrets, Rule{Strip: "vault_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["secret"] != "val" {
		t.Errorf("expected secret=val, got %v", result)
	}
}

func TestApply_StripMakesEmpty(t *testing.T) {
	secrets := map[string]string{"prefix": "val"}
	_, err := Apply(secrets, Rule{Strip: "prefix"})
	if err == nil {
		t.Fatal("expected error when strip makes key empty")
	}
}

func TestApply_CombinedRules(t *testing.T) {
	secrets := map[string]string{"raw_key": "value"}
	result, err := Apply(secrets, Rule{Strip: "raw_", Uppercase: true, Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["APP_KEY"] != "value" {
		t.Errorf("expected APP_KEY=value, got %v", result)
	}
}

func TestApply_EmptyKey(t *testing.T) {
	secrets := map[string]string{"": "val"}
	_, err := Apply(secrets, Rule{})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}
