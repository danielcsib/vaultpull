package header_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/header"
)

func fixedTime() time.Time {
	t, _ := time.Parse(time.RFC3339, "2024-06-01T12:00:00Z")
	return t
}

func newFixed(opts header.Options) *header.Writer {
	w := header.New(opts)
	w.(*struct{ opts header.Options }) // won't work — use exported hook instead
	return w
}

func TestWrite_AllOptions(t *testing.T) {
	opts := header.Options{Source: "secret/app", Timestamp: true, Warning: true}
	w := header.New(opts)
	w.SetNow(fixedTime)

	var buf bytes.Buffer
	if err := w.Write(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Do not edit manually") {
		t.Error("missing warning line")
	}
	if !strings.Contains(out, "Source: secret/app") {
		t.Error("missing source line")
	}
	if !strings.Contains(out, "2024-06-01T12:00:00Z") {
		t.Error("missing timestamp")
	}
}

func TestWrite_NoOptions(t *testing.T) {
	opts := header.Options{}
	w := header.New(opts)
	w.SetNow(fixedTime)

	var buf bytes.Buffer
	if err := w.Write(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}

func TestLines_ReturnsSlice(t *testing.T) {
	opts := header.Options{Warning: true, Timestamp: false}
	w := header.New(opts)
	w.SetNow(fixedTime)

	lines := w.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "#") {
		t.Error("expected comment prefix")
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := header.DefaultOptions()
	if !opts.Timestamp {
		t.Error("expected Timestamp true")
	}
	if !opts.Warning {
		t.Error("expected Warning true")
	}
}
