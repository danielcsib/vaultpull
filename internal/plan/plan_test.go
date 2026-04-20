package plan_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/plan"
)

func TestBuild_AllUnchanged(t *testing.T) {
	current := map[string]string{"A": "1", "B": "2"}
	incoming := map[string]string{"A": "1", "B": "2"}
	p := plan.Build(".env", current, incoming)
	if p.HasChanges() {
		t.Fatal("expected no changes")
	}
	for _, e := range p.Entries {
		if e.Action != plan.ActionUnchanged {
			t.Errorf("key %s: want unchanged, got %s", e.Key, e.Action)
		}
	}
}

func TestBuild_Added(t *testing.T) {
	p := plan.Build(".env", map[string]string{}, map[string]string{"NEW": "val"})
	if len(p.Entries) != 1 || p.Entries[0].Action != plan.ActionAdd {
		t.Fatalf("expected one add entry, got %+v", p.Entries)
	}
	if !p.HasChanges() {
		t.Fatal("expected HasChanges to be true")
	}
}

func TestBuild_Removed(t *testing.T) {
	p := plan.Build(".env", map[string]string{"OLD": "v"}, map[string]string{})
	if len(p.Entries) != 1 || p.Entries[0].Action != plan.ActionDelete {
		t.Fatalf("expected one delete entry, got %+v", p.Entries)
	}
}

func TestBuild_Updated(t *testing.T) {
	p := plan.Build(".env", map[string]string{"K": "old"}, map[string]string{"K": "new"})
	if len(p.Entries) != 1 || p.Entries[0].Action != plan.ActionUpdate {
		t.Fatalf("expected one update entry, got %+v", p.Entries)
	}
	if p.Entries[0].OldVal != "old" || p.Entries[0].NewVal != "new" {
		t.Errorf("unexpected values: %+v", p.Entries[0])
	}
}

func TestBuild_SortedKeys(t *testing.T) {
	incoming := map[string]string{"Z": "1", "A": "2", "M": "3"}
	p := plan.Build(".env", map[string]string{}, incoming)
	keys := make([]string, len(p.Entries))
	for i, e := range p.Entries {
		keys[i] = e.Key
	}
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("keys not sorted: %v", keys)
		}
	}
}

func TestPrint_ContainsSymbols(t *testing.T) {
	current := map[string]string{"EXISTING": "v", "REMOVED": "x"}
	incoming := map[string]string{"EXISTING": "changed", "ADDED": "new"}
	p := plan.Build(".env", current, incoming)
	var buf bytes.Buffer
	p.Print(&buf)
	out := buf.String()
	if !strings.Contains(out, "+ ADDED") {
		t.Errorf("missing add symbol in output:\n%s", out)
	}
	if !strings.Contains(out, "~ EXISTING") {
		t.Errorf("missing update symbol in output:\n%s", out)
	}
	if !strings.Contains(out, "- REMOVED") {
		t.Errorf("missing delete symbol in output:\n%s", out)
	}
}

func TestPrint_EmptyPlan(t *testing.T) {
	p := plan.Build(".env", nil, nil)
	var buf bytes.Buffer
	p.Print(&buf)
	if !strings.Contains(buf.String(), "no keys") {
		t.Errorf("expected '(no keys)' in output: %s", buf.String())
	}
}
