// Package plan computes a dry-run plan of changes before secrets are written.
package plan

import (
	"fmt"
	"io"
	"sort"
)

// Action describes what will happen to a secret key.
type Action string

const (
	ActionAdd    Action = "add"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
	ActionUnchanged Action = "unchanged"
)

// Entry represents a single planned change.
type Entry struct {
	Key    string
	Action Action
	OldVal string // masked or empty
	NewVal string // masked or empty
}

// Plan holds all planned changes for a single env file path.
type Plan struct {
	Path    string
	Entries []Entry
}

// Build compares current (existing) secrets with incoming (new) secrets
// and returns a Plan describing what would change.
func Build(path string, current, incoming map[string]string) Plan {
	p := Plan{Path: path}

	for k, newVal := range incoming {
		if oldVal, ok := current[k]; !ok {
			p.Entries = append(p.Entries, Entry{Key: k, Action: ActionAdd, NewVal: newVal})
		} else if oldVal != newVal {
			p.Entries = append(p.Entries, Entry{Key: k, Action: ActionUpdate, OldVal: oldVal, NewVal: newVal})
		} else {
			p.Entries = append(p.Entries, Entry{Key: k, Action: ActionUnchanged})
		}
	}

	for k := range current {
		if _, ok := incoming[k]; !ok {
			p.Entries = append(p.Entries, Entry{Key: k, Action: ActionDelete, OldVal: current[k]})
		}
	}

	sort.Slice(p.Entries, func(i, j int) bool {
		return p.Entries[i].Key < p.Entries[j].Key
	})
	return p
}

// HasChanges returns true if any entry is not ActionUnchanged.
func (p Plan) HasChanges() bool {
	for _, e := range p.Entries {
		if e.Action != ActionUnchanged {
			return true
		}
	}
	return false
}

// Print writes a human-readable summary of the plan to w.
func (p Plan) Print(w io.Writer) {
	fmt.Fprintf(w, "Plan for %s:\n", p.Path)
	if len(p.Entries) == 0 {
		fmt.Fprintln(w, "  (no keys)")
		return
	}
	for _, e := range p.Entries {
		switch e.Action {
		case ActionAdd:
			fmt.Fprintf(w, "  + %s\n", e.Key)
		case ActionUpdate:
			fmt.Fprintf(w, "  ~ %s\n", e.Key)
		case ActionDelete:
			fmt.Fprintf(w, "  - %s\n", e.Key)
		case ActionUnchanged:
			fmt.Fprintf(w, "    %s\n", e.Key)
		}
	}
}
