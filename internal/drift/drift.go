// Package drift detects when local .env files have diverged from Vault secrets.
package drift

import (
	"fmt"
	"sort"
)

// Status represents the drift state of a single key.
type Status int

const (
	StatusMatch   Status = iota // value matches Vault
	StatusDrifted               // local value differs from Vault
	StatusMissing               // key absent in local file
	StatusOrphan                // key present locally but not in Vault
)

func (s Status) String() string {
	switch s {
	case StatusMatch:
		return "match"
	case StatusDrifted:
		return "drifted"
	case StatusMissing:
		return "missing"
	case StatusOrphan:
		return "orphan"
	default:
		return "unknown"
	}
}

// Entry describes the drift state for one key.
type Entry struct {
	Key    string
	Status Status
}

// Report holds the full drift analysis for a path.
type Report struct {
	Path    string
	Entries []Entry
}

// HasDrift returns true if any entry is not StatusMatch.
func (r *Report) HasDrift() bool {
	for _, e := range r.Entries {
		if e.Status != StatusMatch {
			return true
		}
	}
	return false
}

// Drifted returns only entries whose status is not StatusMatch.
func (r *Report) Drifted() []Entry {
	out := make([]Entry, 0)
	for _, e := range r.Entries {
		if e.Status != StatusMatch {
			out = append(out, e)
		}
	}
	return out
}

// Detect compares vault secrets against local env values and returns a Report.
// vault is the source of truth; local is what is currently on disk.
func Detect(path string, vault, local map[string]string) (*Report, error) {
	if path == "" {
		return nil, fmt.Errorf("drift: path must not be empty")
	}

	seen := make(map[string]bool)
	var entries []Entry

	for k, vv := range vault {
		seen[k] = true
		lv, ok := local[k]
		if !ok {
			entries = append(entries, Entry{Key: k, Status: StatusMissing})
		} else if lv != vv {
			entries = append(entries, Entry{Key: k, Status: StatusDrifted})
		} else {
			entries = append(entries, Entry{Key: k, Status: StatusMatch})
		}
	}

	for k := range local {
		if !seen[k] {
			entries = append(entries, Entry{Key: k, Status: StatusOrphan})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return &Report{Path: path, Entries: entries}, nil
}
