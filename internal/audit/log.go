package audit

import (
	"encoding/json"
	"os"
	"time"
)

// Entry represents a single audit log entry for a sync operation.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	SecretPath string    `json:"secret_path"`
	EnvFile    string    `json:"env_file"`
	Keys       []string  `json:"keys"`
	Success    bool      `json:"success"`
	Error      string    `json:"error,omitempty"`
}

// Logger writes audit entries to a file in JSON Lines format.
type Logger struct {
	path string
}

// NewLogger creates a Logger that appends to the given file path.
func NewLogger(path string) *Logger {
	return &Logger{path: path}
}

// Record appends an audit entry to the log file.
func (l *Logger) Record(e Entry) error {
	if l.path == "" {
		return nil
	}
	e.Timestamp = time.Now().UTC()

	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	return enc.Encode(e)
}
