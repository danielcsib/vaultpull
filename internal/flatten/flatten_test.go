package flatten_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/flatten"
)

func TestFlatten_EmptyMap(t *testing.T) {
	f := flatten.New(flatten.DefaultOptions())
	out := f.Flatten(map[string]any{})
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %v", out)
	}
}

func TestFlatten_FlatInput(t *testing.T) {
	f := flatten.New(flatten.DefaultOptions())
	out := f.Flatten(map[string]any{
		"host": "localhost",
		"port": "5432",
	})
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", out["HOST"])
	}
	if out["PORT"] != "5432" {
		t.Errorf("expected PORT=5432, got %q", out["PORT"])
	}
}

func TestFlatten_NestedMap(t *testing.T) {
	f := flatten.New(flatten.DefaultOptions())
	out := f.Flatten(map[string]any{
		"db": map[string]any{
			"host": "db.internal",
			"port": "5432",
		},
	})
	if out["DB_HOST"] != "db.internal" {
		t.Errorf("expected DB_HOST=db.internal, got %q", out["DB_HOST"])
	}
	if out["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", out["DB_PORT"])
	}
}

func TestFlatten_DeeplyNested(t *testing.T) {
	f := flatten.New(flatten.DefaultOptions())
	out := f.Flatten(map[string]any{
		"app": map[string]any{
			"cache": map[string]any{
				"ttl": "300",
			},
		},
	})
	if out["APP_CACHE_TTL"] != "300" {
		t.Errorf("expected APP_CACHE_TTL=300, got %q", out["APP_CACHE_TTL"])
	}
}

func TestFlatten_StringStringMap(t *testing.T) {
	f := flatten.New(flatten.DefaultOptions())
	out := f.Flatten(map[string]any{
		"labels": map[string]string{
			"env": "prod",
		},
	})
	if out["LABELS_ENV"] != "prod" {
		t.Errorf("expected LABELS_ENV=prod, got %q", out["LABELS_ENV"])
	}
}

func TestFlatten_CustomSeparator(t *testing.T) {
	f := flatten.New(flatten.Options{Separator: ".", UppercaseKeys: false})
	out := f.Flatten(map[string]any{
		"db": map[string]any{
			"host": "localhost",
		},
	})
	if out["db.host"] != "localhost" {
		t.Errorf("expected db.host=localhost, got %v", out)
	}
}

func TestFlatten_LowercaseKeys(t *testing.T) {
	f := flatten.New(flatten.Options{Separator: "_", UppercaseKeys: false})
	out := f.Flatten(map[string]any{"Key": "val"})
	if out["Key"] != "val" {
		t.Errorf("expected Key=val, got %v", out)
	}
}

func TestFlatten_NonStringValue(t *testing.T) {
	f := flatten.New(flatten.DefaultOptions())
	out := f.Flatten(map[string]any{"timeout": 30})
	if out["TIMEOUT"] != "30" {
		t.Errorf("expected TIMEOUT=30, got %q", out["TIMEOUT"])
	}
}
