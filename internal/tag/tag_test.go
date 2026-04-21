package tag_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/tag"
)

func TestSet_And_Get(t *testing.T) {
	tr := tag.New()
	tr.Set("DB_PASSWORD", "env", "prod")
	tr.Set("DB_PASSWORD", "team", "platform")

	tags := tr.Get("DB_PASSWORD")
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
}

func TestSet_OverwritesDuplicateTagKey(t *testing.T) {
	tr := tag.New()
	tr.Set("API_KEY", "env", "staging")
	tr.Set("API_KEY", "env", "prod")

	tags := tr.Get("API_KEY")
	if len(tags) != 1 {
		t.Fatalf("expected 1 tag after overwrite, got %d", len(tags))
	}
	if tags[0].Value != "prod" {
		t.Errorf("expected value 'prod', got %q", tags[0].Value)
	}
}

func TestSet_EmptyKeyIsNoOp(t *testing.T) {
	tr := tag.New()
	tr.Set("", "env", "prod")
	tr.Set("DB_PASSWORD", "", "prod")

	if len(tr.Get("")) != 0 {
		t.Error("expected no tags for empty secret key")
	}
	if len(tr.Get("DB_PASSWORD")) != 0 {
		t.Error("expected no tags when tag key is empty")
	}
}

func TestHas_ReturnsTrueOnMatch(t *testing.T) {
	tr := tag.New()
	tr.Set("TOKEN", "sensitivity", "high")

	if !tr.Has("TOKEN", "sensitivity", "high") {
		t.Error("expected Has to return true")
	}
	if tr.Has("TOKEN", "sensitivity", "low") {
		t.Error("expected Has to return false for wrong value")
	}
}

func TestFilter_ReturnsMatchingKeys(t *testing.T) {
	tr := tag.New()
	tr.Set("DB_PASS", "env", "prod")
	tr.Set("DB_PASS", "team", "platform")
	tr.Set("API_KEY", "env", "prod")
	tr.Set("DEBUG", "env", "dev")

	m := map[string]string{
		"DB_PASS": "secret1",
		"API_KEY": "secret2",
		"DEBUG":   "true",
	}

	result := tr.Filter(m, []tag.Tag{{Key: "env", Value: "prod"}})
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
	if _, ok := result["DEBUG"]; ok {
		t.Error("DEBUG should not be in filtered result")
	}
}

func TestFilter_ANDSemantics(t *testing.T) {
	tr := tag.New()
	tr.Set("DB_PASS", "env", "prod")
	tr.Set("DB_PASS", "team", "platform")
	tr.Set("API_KEY", "env", "prod")

	m := map[string]string{"DB_PASS": "s1", "API_KEY": "s2"}

	constraints := []tag.Tag{
		{Key: "env", Value: "prod"},
		{Key: "team", Value: "platform"},
	}
	result := tr.Filter(m, constraints)
	if len(result) != 1 {
		t.Fatalf("expected 1 result with AND constraints, got %d", len(result))
	}
	if _, ok := result["DB_PASS"]; !ok {
		t.Error("expected DB_PASS in result")
	}
}

func TestString_NoTags(t *testing.T) {
	tr := tag.New()
	s := tr.String("UNKNOWN")
	if s != "UNKNOWN: (no tags)" {
		t.Errorf("unexpected string: %q", s)
	}
}
