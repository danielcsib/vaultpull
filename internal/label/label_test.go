package label_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/label"
)

func TestApply_InjectsLabels(t *testing.T) {
	l := label.New(label.Set{"ENV": "production", "TEAM": "platform"}, "")
	secrets := map[string]string{"DB_PASS": "secret"}
	out := l.Apply(secrets)

	if out["DB_PASS"] != "secret" {
		t.Errorf("original key missing")
	}
	if out["VAULTPULL_ENV"] != "production" {
		t.Errorf("expected VAULTPULL_ENV=production, got %q", out["VAULTPULL_ENV"])
	}
	if out["VAULTPULL_TEAM"] != "platform" {
		t.Errorf("expected VAULTPULL_TEAM=platform")
	}
}

func TestApply_DoesNotOverwriteExisting(t *testing.T) {
	l := label.New(label.Set{"ENV": "staging"}, "")
	secrets := map[string]string{"VAULTPULL_ENV": "already-set"}
	out := l.Apply(secrets)
	if out["VAULTPULL_ENV"] != "already-set" {
		t.Errorf("existing key should not be overwritten")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	l := label.New(label.Set{"X": "1"}, "")
	secrets := map[string]string{"A": "b"}
	l.Apply(secrets)
	if _, ok := secrets["VAULTPULL_X"]; ok {
		t.Errorf("input map was mutated")
	}
}

func TestStrip_RemovesLabelKeys(t *testing.T) {
	l := label.New(label.Set{}, "")
	secrets := map[string]string{
		"DB_PASS":       "s3cr3t",
		"VAULTPULL_ENV": "prod",
	}
	out := l.Strip(secrets)
	if _, ok := out["VAULTPULL_ENV"]; ok {
		t.Errorf("label key should have been stripped")
	}
	if out["DB_PASS"] != "s3cr3t" {
		t.Errorf("non-label key should remain")
	}
}

func TestNew_CustomPrefix(t *testing.T) {
	l := label.New(label.Set{"VER": "2"}, "META_")
	out := l.Apply(map[string]string{})
	if out["META_VER"] != "2" {
		t.Errorf("expected META_VER=2, got %q", out["META_VER"])
	}
}
