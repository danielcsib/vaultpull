package env_test

import (
	"os"
	"testing"

	"github.com/yourusername/vaultpull/internal/env"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_EmptyPath(t *testing.T) {
	_, err := env.Load("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	got, err := env.Load("/nonexistent/path/.env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty map, got %v", got)
	}
}

func TestLoad_BasicKeyValues(t *testing.T) {
	path := writeTemp(t, "FOO=bar\nBAZ=qux\n")
	got, err := env.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if got["FOO"] != "bar" || got["BAZ"] != "qux" {
		t.Fatalf("unexpected map: %v", got)
	}
}

func TestLoad_IgnoresComments(t *testing.T) {
	path := writeTemp(t, "# comment\nKEY=value\n")
	got, err := env.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got["KEY"] != "value" {
		t.Fatalf("unexpected map: %v", got)
	}
}

func TestLoad_QuotedValues(t *testing.T) {
	path := writeTemp(t, `A="hello world"
B='single quoted'
`)
	got, err := env.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if got["A"] != "hello world" {
		t.Errorf("A: got %q", got["A"])
	}
	if got["B"] != "single quoted" {
		t.Errorf("B: got %q", got["B"])
	}
}

func TestLoad_InvalidLine(t *testing.T) {
	path := writeTemp(t, "INVALID\n")
	_, err := env.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid line")
	}
}
