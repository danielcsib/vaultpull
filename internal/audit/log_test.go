package audit_test

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"

	"github.com/your-org/vaultpull/internal/audit"
)

func TestRecord_WritesEntry(t *testing.T) {
	tmp, err := os.CreateTemp("", "audit-*.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	l := audit.NewLogger(tmp.Name())
	err = l.Record(audit.Entry{
		SecretPath: "secret/app",
		EnvFile:    ".env",
		Keys:       []string{"DB_PASS", "API_KEY"},
		Success:    true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	f, _ := os.Open(tmp.Name())
	defer f.Close()
	var entry audit.Entry
	if err := json.NewDecoder(bufio.NewReader(f)).Decode(&entry); err != nil {
		t.Fatalf("failed to decode entry: %v", err)
	}
	if entry.SecretPath != "secret/app" {
		t.Errorf("expected secret/app, got %s", entry.SecretPath)
	}
	if !entry.Success {
		t.Error("expected success=true")
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestRecord_EmptyPath_NoOp(t *testing.T) {
	l := audit.NewLogger("")
	if err := l.Record(audit.Entry{SecretPath: "x", Success: true}); err != nil {
		t.Errorf("expected no error for empty path, got %v", err)
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	tmp, _ := os.CreateTemp("", "audit-multi-*.jsonl")
	tmp.Close()
	defer os.Remove(tmp.Name())

	l := audit.NewLogger(tmp.Name())
	for i := 0; i < 3; i++ {
		if err := l.Record(audit.Entry{SecretPath: "s", Success: true}); err != nil {
			t.Fatal(err)
		}
	}

	f, _ := os.Open(tmp.Name())
	defer f.Close()
	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		count++
	}
	if count != 3 {
		t.Errorf("expected 3 lines, got %d", count)
	}
}
