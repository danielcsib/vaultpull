package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/output"
)

func TestWrite_EnvFormat(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(output.FormatEnv, &buf)
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost in output, got:\n%s", got)
	}
	if !strings.Contains(got, "DB_PORT=5432") {
		t.Errorf("expected DB_PORT=5432 in output, got:\n%s", got)
	}
}

func TestWrite_ExportFormat(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(output.FormatExport, &buf)
	if err := w.Write(map[string]string{"API_KEY": "abc"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(buf.String(), "export ") {
		t.Errorf("expected export prefix, got: %s", buf.String())
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(output.FormatJSON, &buf)
	if err := w.Write(map[string]string{"TOKEN": "secret"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "\"TOKEN\"") {
		t.Errorf("expected TOKEN key in JSON, got:\n%s", got)
	}
	if !strings.Contains(got, "\"secret\"") {
		t.Errorf("expected secret value in JSON, got:\n%s", got)
	}
}

func TestWrite_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(output.Format("xml"), &buf)
	if err := w.Write(map[string]string{"K": "v"}); err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}

func TestWrite_QuotedValues(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(output.FormatEnv, &buf)
	if err := w.Write(map[string]string{"MSG": "hello world"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"hello world"`) {
		t.Errorf("expected quoted value, got: %s", buf.String())
	}
}

func TestNew_NilWriterUsesStdout(t *testing.T) {
	w := output.New(output.FormatEnv, nil)
	if w == nil {
		t.Fatal("expected non-nil Writer")
	}
}
