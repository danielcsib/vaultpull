package scope_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/scope"
)

func TestNew_EmptyRoot(t *testing.T) {
	_, err := scope.New("")
	if err == nil {
		t.Fatal("expected error for empty root, got nil")
	}
}

func TestNew_WhitespaceRoot(t *testing.T) {
	_, err := scope.New("   ")
	if err == nil {
		t.Fatal("expected error for whitespace root, got nil")
	}
}

func TestNew_ValidRoot(t *testing.T) {
	s, err := scope.New("secret/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Root() != "secret/myapp" {
		t.Errorf("expected root %q, got %q", "secret/myapp", s.Root())
	}
}

func TestNew_TrimsTrailingSlash(t *testing.T) {
	s, err := scope.New("secret/myapp/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Root() != "secret/myapp" {
		t.Errorf("expected trimmed root, got %q", s.Root())
	}
}

func TestContains_ExactMatch(t *testing.T) {
	s, _ := scope.New("secret/myapp")
	if !s.Contains("secret/myapp") {
		t.Error("expected exact root to be contained")
	}
}

func TestContains_ChildPath(t *testing.T) {
	s, _ := scope.New("secret/myapp")
	if !s.Contains("secret/myapp/db") {
		t.Error("expected child path to be contained")
	}
}

func TestContains_SiblingPath(t *testing.T) {
	s, _ := scope.New("secret/myapp")
	if s.Contains("secret/myapp2") {
		t.Error("sibling path should not be contained")
	}
}

func TestContains_UnrelatedPath(t *testing.T) {
	s, _ := scope.New("secret/myapp")
	if s.Contains("secret/other") {
		t.Error("unrelated path should not be contained")
	}
}

func TestFilter_ReturnsInScopePaths(t *testing.T) {
	s, _ := scope.New("secret/myapp")
	paths := []string{
		"secret/myapp/db",
		"secret/other/db",
		"secret/myapp/cache",
		"secret/unrelated",
	}
	got := s.Filter(paths)
	if len(got) != 2 {
		t.Fatalf("expected 2 paths, got %d: %v", len(got), got)
	}
}

func TestRelativePath_StripsRoot(t *testing.T) {
	s, _ := scope.New("secret/myapp")
	got := s.RelativePath("secret/myapp/db/credentials")
	if got != "db/credentials" {
		t.Errorf("expected %q, got %q", "db/credentials", got)
	}
}

func TestRelativePath_ExactRoot_ReturnsEmpty(t *testing.T) {
	s, _ := scope.New("secret/myapp")
	got := s.RelativePath("secret/myapp")
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestRelativePath_OutOfScope_Passthrough(t *testing.T) {
	s, _ := scope.New("secret/myapp")
	got := s.RelativePath("secret/other/path")
	if got != "secret/other/path" {
		t.Errorf("expected passthrough, got %q", got)
	}
}
