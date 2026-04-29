package drift_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/drift"
)

func TestDetect_EmptyPath(t *testing.T) {
	_, err := drift.Detect("", map[string]string{"K": "v"}, map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestDetect_AllMatch(t *testing.T) {
	vault := map[string]string{"A": "1", "B": "2"}
	local := map[string]string{"A": "1", "B": "2"}

	r, err := drift.Detect("secrets/app", vault, local)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.HasDrift() {
		t.Error("expected no drift")
	}
	if len(r.Drifted()) != 0 {
		t.Errorf("expected 0 drifted entries, got %d", len(r.Drifted()))
	}
}

func TestDetect_DriftedValue(t *testing.T) {
	vault := map[string]string{"DB_PASS": "secret"}
	local := map[string]string{"DB_PASS": "old"}

	r, err := drift.Detect("secrets/db", vault, local)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.HasDrift() {
		t.Fatal("expected drift")
	}
	if len(r.Drifted()) != 1 || r.Drifted()[0].Status != drift.StatusDrifted {
		t.Errorf("expected one drifted entry, got %+v", r.Drifted())
	}
}

func TestDetect_MissingLocal(t *testing.T) {
	vault := map[string]string{"NEW_KEY": "val"}
	local := map[string]string{}

	r, err := drift.Detect("secrets/app", vault, local)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	d := r.Drifted()
	if len(d) != 1 || d[0].Status != drift.StatusMissing {
		t.Errorf("expected one missing entry, got %+v", d)
	}
}

func TestDetect_OrphanLocal(t *testing.T) {
	vault := map[string]string{}
	local := map[string]string{"LEGACY": "x"}

	r, err := drift.Detect("secrets/app", vault, local)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	d := r.Drifted()
	if len(d) != 1 || d[0].Status != drift.StatusOrphan {
		t.Errorf("expected one orphan entry, got %+v", d)
	}
}

func TestDetect_SortedOutput(t *testing.T) {
	vault := map[string]string{"Z": "1", "A": "2", "M": "3"}
	local := map[string]string{"Z": "1", "A": "2", "M": "3"}

	r, err := drift.Detect("secrets/app", vault, local)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := make([]string, len(r.Entries))
	for i, e := range r.Entries {
		keys[i] = e.Key
	}
	expected := []string{"A", "M", "Z"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("expected key[%d]=%q, got %q", i, k, keys[i])
		}
	}
}

func TestStatus_String(t *testing.T) {
	cases := []struct {
		s    drift.Status
		want string
	}{
		{drift.StatusMatch, "match"},
		{drift.StatusDrifted, "drifted"},
		{drift.StatusMissing, "missing"},
		{drift.StatusOrphan, "orphan"},
	}
	for _, tc := range cases {
		if got := tc.s.String(); got != tc.want {
			t.Errorf("Status(%d).String() = %q, want %q", tc.s, got, tc.want)
		}
	}
}
