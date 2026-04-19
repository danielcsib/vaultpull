// Package watch polls Vault paths at a configurable interval and triggers
// a callback when secrets change relative to a snapshot.
package watch

import (
	"context"
	"log"
	"time"
)

// SecretReader fetches secrets from a Vault path.
type SecretReader interface {
	ReadSecrets(path string) (map[string]string, error)
}

// OnChange is called with the path and new secrets when a change is detected.
type OnChange func(path string, secrets map[string]string)

// Watcher polls one or more Vault paths for changes.
type Watcher struct {
	client   SecretReader
	paths    []string
	interval time.Duration
	onChange OnChange
	last     map[string]map[string]string
}

// New creates a Watcher. interval must be > 0.
func New(client SecretReader, paths []string, interval time.Duration, fn OnChange) *Watcher {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	return &Watcher{
		client:   client,
		paths:    paths,
		interval: interval,
		onChange: fn,
		last:     make(map[string]map[string]string),
	}
}

// Run starts polling until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	w.poll()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.poll()
		}
	}
}

func (w *Watcher) poll() {
	for _, path := range w.paths {
		secrets, err := w.client.ReadSecrets(path)
		if err != nil {
			log.Printf("watch: error reading %s: %v", path, err)
			continue
		}
		if changed(w.last[path], secrets) {
			w.last[path] = copyMap(secrets)
			if w.onChange != nil {
				w.onChange(path, secrets)
			}
		}
	}
}

func changed(prev, next map[string]string) bool {
	if len(prev) != len(next) {
		return true
	}
	for k, v := range next {
		if prev[k] != v {
			return true
		}
	}
	return false
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
