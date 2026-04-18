package merge_test

import (
	"os"
	"testing"

	"github.com/your-org/vaultpull/internal/merge"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestFromFile_NoExistingFile(t *testing.T) {
	incoming := map[string]string{"FOO": "bar"}
	res, err := merge.FromFile("/nonexistent/.env", incoming)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Final["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", res.Final["FOO"])
	}
	if len(res.Added) != 1 || res.Added[0] != "FOO" {
		t.Errorf("expected FOO in Added, got %v", res.Added)
	}
}

func TestFromFile_PreservesLocalKeys(t *testing.T) {
	path := writeTemp(t, "LOCAL_ONLY=secret\nSHARED=old\n")
	incoming := map[string]string{"SHARED": "new"}
	res, err := merge.FromFile(path, incoming)
	if err != nil {
		t.Fatal(err)
	}
	if res.Final["LOCAL_ONLY"] != "secret" {
		t.Error("LOCAL_ONLY should be preserved")
	}
	if len(res.Preserved) != 1 || res.Preserved[0] != "LOCAL_ONLY" {
		t.Errorf("expected LOCAL_ONLY in Preserved, got %v", res.Preserved)
	}
}

func TestFromFile_DetectsUpdated(t *testing.T) {
	path := writeTemp(t, "KEY=oldval\n")
	incoming := map[string]string{"KEY": "newval"}
	res, err := merge.FromFile(path, incoming)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Updated) != 1 || res.Updated[0] != "KEY" {
		t.Errorf("expected KEY in Updated, got %v", res.Updated)
	}
	if res.Final["KEY"] != "newval" {
		t.Errorf("expected newval, got %q", res.Final["KEY"])
	}
}

func TestFromFile_QuotedValues(t *testing.T) {
	path := writeTemp(t, `KEY="quoted value"`+"\n")
	incoming := map[string]string{}
	res, err := merge.FromFile(path, incoming)
	if err != nil {
		t.Fatal(err)
	}
	if res.Final["KEY"] != "quoted value" {
		t.Errorf("expected unquoted value, got %q", res.Final["KEY"])
	}
}

func TestFromFile_IgnoresComments(t *testing.T) {
	path := writeTemp(t, "# comment\nVALID=yes\n")
	res, err := merge.FromFile(path, map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := res.Final["# comment"]; ok {
		t.Error("comment line should not be parsed as key")
	}
	if res.Final["VALID"] != "yes" {
		t.Error("VALID key should be present")
	}
}
