package envwriter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteEnvFile_BasicSecrets(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_NAME": "vaultpull",
	}

	if err := WriteEnvFile(path, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}

	for key, val := range secrets {
		expected := key + "=" + val
		if !strings.Contains(string(content), expected) {
			t.Errorf("expected %q in output, got:\n%s", expected, content)
		}
	}
}

func TestWriteEnvFile_QuotedValues(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	secrets := map[string]string{
		"SECRET_KEY": "hello world",
		"MULTILINE":  "line1\nline2",
	}

	if err := WriteEnvFile(path, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := string(must(os.ReadFile(path)))
	if !strings.Contains(content, `SECRET_KEY="hello world"`) {
		t.Errorf("expected quoted value for SECRET_KEY, got:\n%s", content)
	}
}

func TestWriteEnvFile_EmptyPath(t *testing.T) {
	err := WriteEnvFile("", map[string]string{"K": "V"})
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestWriteEnvFile_SortedOutput(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	secrets := map[string]string{
		"ZEBRA": "last",
		"ALPHA": "first",
		"MANGO": "middle",
	}

	if err := WriteEnvFile(path, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := string(must(os.ReadFile(path)))
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) != 3 || !strings.HasPrefix(lines[0], "ALPHA") || !strings.HasPrefix(lines[2], "ZEBRA") {
		t.Errorf("expected sorted output, got:\n%s", content)
	}
}

func must(b []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return b
}
