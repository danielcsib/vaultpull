package schema_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/schema"
)

func TestValidate_AllPresent(t *testing.T) {
	s := schema.New([]schema.FieldRule{
		{Key: "DB_HOST", Required: true},
		{Key: "DB_PORT", Required: true},
	})
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	violations, err := s.Validate(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestValidate_MissingRequired(t *testing.T) {
	s := schema.New([]schema.FieldRule{
		{Key: "API_KEY", Required: true},
	})
	violations, err := s.Validate(map[string]string{})
	if err == nil {
		t.Fatal("expected error")
	}
	if len(violations) != 1 || violations[0].Key != "API_KEY" {
		t.Fatalf("unexpected violations: %v", violations)
	}
}

func TestValidate_EmptyValueDisallowed(t *testing.T) {
	s := schema.New([]schema.FieldRule{
		{Key: "TOKEN", Required: true, AllowEmpty: false},
	})
	violations, err := s.Validate(map[string]string{"TOKEN": "   "})
	if err == nil {
		t.Fatal("expected error")
	}
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidate_EmptyValueAllowed(t *testing.T) {
	s := schema.New([]schema.FieldRule{
		{Key: "OPTIONAL", Required: false, AllowEmpty: true},
	})
	violations, err := s.Validate(map[string]string{"OPTIONAL": ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("expected no violations")
	}
}

func TestValidate_MultipleViolations(t *testing.T) {
	s := schema.New([]schema.FieldRule{
		{Key: "A", Required: true},
		{Key: "B", Required: true},
	})
	violations, err := s.Validate(map[string]string{})
	if err == nil {
		t.Fatal("expected error")
	}
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(violations))
	}
}
