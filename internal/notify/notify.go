package notify

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Event holds details of a sync notification.
type Event struct {
	Level   Level
	Path    string
	Message string
	Time    time.Time
}

// Notifier writes sync events to an output sink.
type Notifier struct {
	out io.Writer
}

// New returns a Notifier writing to w. If w is nil, os.Stdout is used.
func New(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{out: w}
}

// Send emits an event to the configured writer.
func (n *Notifier) Send(e Event) {
	if e.Time.IsZero() {
		e.Time = time.Now()
	}
	timestamp := e.Time.UTC().Format(time.RFC3339)
	path := e.Path
	if path == "" {
		path = "-"
	}
	fmt.Fprintf(n.out, "%s [%s] path=%s %s\n", timestamp, strings.ToUpper(string(e.Level)), path, e.Message)
}

// Info is a convenience wrapper for LevelInfo events.
func (n *Notifier) Info(path, msg string) {
	n.Send(Event{Level: LevelInfo, Path: path, Message: msg})
}

// Warn is a convenience wrapper for LevelWarn events.
func (n *Notifier) Warn(path, msg string) {
	n.Send(Event{Level: LevelWarn, Path: path, Message: msg})
}

// Error is a convenience wrapper for LevelError events.
func (n *Notifier) Error(path, msg string) {
	n.Send(Event{Level: LevelError, Path: path, Message: msg})
}
