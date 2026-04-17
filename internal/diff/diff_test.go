package diff

import (
	"testing"
)

func TestCompare_AllUnchanged(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	next := map[string]string{"A": "1", "B": "2"}
	r := Compare(old, next)
	if r.HasChanges() {
		t.Fatal("expected no changes")
	}
	if len(r.Unchanged) != 2 {
		t.Fatalf("expected 2 unchanged, got %d", len(r.Unchanged))
	}
}

func TestCompare_Added(t *testing.T) {
	old := map[string]string{}
	next := map[string]string{"NEW_KEY": "val"}
	r := Compare(old, next)
	if len(r.Added) != 1 || r.Added[0] != "NEW_KEY" {
		t.Fatalf("expected NEW_KEY added, got %v", r.Added)
	}
}

func TestCompare_Removed(t *testing.T) {
	old := map[string]string{"OLD": "x"}
	next := map[string]string{}
	r := Compare(old, next)
	if len(r.Removed) != 1 || r.Removed[0] != "OLD" {
		t.Fatalf("expected OLD removed, got %v", r.Removed)
	}
}

func TestCompare_Changed(t *testing.T) {
	old := map[string]string{"K": "old_val"}
	next := map[string]string{"K": "new_val"}
	r := Compare(old, next)
	if len(r.Changed) != 1 || r.Changed[0] != "K" {
		t.Fatalf("expected K changed, got %v", r.Changed)
	}
	if !r.HasChanges() {
		t.Fatal("expected HasChanges to be true")
	}
}

func TestCompare_Mixed(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2", "C": "3"}
	next := map[string]string{"A": "1", "B": "updated", "D": "4"}
	r := Compare(old, next)
	if len(r.Unchanged) != 1 || r.Unchanged[0] != "A" {
		t.Fatalf("unexpected unchanged: %v", r.Unchanged)
	}
	if len(r.Changed) != 1 || r.Changed[0] != "B" {
		t.Fatalf("unexpected changed: %v", r.Changed)
	}
	if len(r.Added) != 1 || r.Added[0] != "D" {
		t.Fatalf("unexpected added: %v", r.Added)
	}
	if len(r.Removed) != 1 || r.Removed[0] != "C" {
		t.Fatalf("unexpected removed: %v", r.Removed)
	}
}

func TestCompare_BothEmpty(t *testing.T) {
	r := Compare(map[string]string{}, map[string]string{})
	if r.HasChanges() {
		t.Fatal("expected no changes for empty maps")
	}
}
