package notify

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	n := New(nil)
	if n.out == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestSend_FormatsOutput(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)
	fixed := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	n.Send(Event{
		Level:   LevelInfo,
		Path:    "secret/app",
		Message: "synced 3 keys",
		Time:    fixed,
	})
	out := buf.String()
	if !strings.Contains(out, "[INFO]") {
		t.Errorf("expected INFO in output, got: %s", out)
	}
	if !strings.Contains(out, "path=secret/app") {
		t.Errorf("expected path in output, got: %s", out)
	}
	if !strings.Contains(out, "synced 3 keys") {
		t.Errorf("expected message in output, got: %s", out)
	}
}

func TestSend_EmptyPath_UsesDash(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)
	n.Send(Event{Level: LevelWarn, Message: "no path", Time: time.Now()})
	if !strings.Contains(buf.String(), "path=-") {
		t.Errorf("expected path=- for empty path, got: %s", buf.String())
	}
}

func TestInfo_SetsLevel(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)
	n.Info("secret/db", "ok")
	if !strings.Contains(buf.String(), "[INFO]") {
		t.Errorf("expected INFO level")
	}
}

func TestWarn_SetsLevel(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)
	n.Warn("secret/db", "stale cache")
	if !strings.Contains(buf.String(), "[WARN]") {
		t.Errorf("expected WARN level")
	}
}

func TestError_SetsLevel(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)
	n.Error("secret/db", "vault unreachable")
	if !strings.Contains(buf.String(), "[ERROR]") {
		t.Errorf("expected ERROR level")
	}
}

func TestSend_AutoTimestamp(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)
	n.Send(Event{Level: LevelInfo, Message: "auto time"})
	if buf.Len() == 0 {
		t.Fatal("expected output")
	}
}
