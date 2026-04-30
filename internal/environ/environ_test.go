package environ_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpull/internal/environ"
)

func TestOverlay_SecretsOverwriteByDefault(t *testing.T) {
	t.Setenv("EXISTING_KEY", "os-value")

	secrets := map[string]string{"EXISTING_KEY": "vault-value", "NEW_KEY": "new"}
	result := environ.Overlay(secrets, environ.DefaultOptions())

	if result["EXISTING_KEY"] != "vault-value" {
		t.Errorf("expected vault-value, got %q", result["EXISTING_KEY"])
	}
	if result["NEW_KEY"] != "new" {
		t.Errorf("expected new, got %q", result["NEW_KEY"])
	}
}

func TestOverlay_OSWinsWhenOverwriteFalse(t *testing.T) {
	t.Setenv("EXISTING_KEY", "os-value")

	secrets := map[string]string{"EXISTING_KEY": "vault-value"}
	opts := environ.Options{Overwrite: false}
	result := environ.Overlay(secrets, opts)

	if result["EXISTING_KEY"] != "os-value" {
		t.Errorf("expected os-value, got %q", result["EXISTING_KEY"])
	}
}

func TestOverlay_PrefixFiltersOSEnv(t *testing.T) {
	t.Setenv("APP_FOO", "1")
	t.Setenv("OTHER_BAR", "2")

	opts := environ.Options{Overwrite: true, Prefix: "APP_"}
	result := environ.Overlay(map[string]string{}, opts)

	if _, ok := result["APP_FOO"]; !ok {
		t.Error("expected APP_FOO to be present")
	}
	if _, ok := result["OTHER_BAR"]; ok {
		t.Error("expected OTHER_BAR to be excluded by prefix filter")
	}
}

func TestOverlay_DoesNotMutateSecrets(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	environ.Overlay(secrets, environ.DefaultOptions())
	if len(secrets) != 1 || secrets["KEY"] != "val" {
		t.Error("Overlay mutated the secrets map")
	}
}

func TestToSlice_ProducesKeyValuePairs(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	slice := environ.ToSlice(env)

	if len(slice) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(slice))
	}
	for _, entry := range slice {
		if !strings.Contains(entry, "=") {
			t.Errorf("entry missing '=': %q", entry)
		}
	}
}

func TestToSlice_EmptyMap(t *testing.T) {
	slice := environ.ToSlice(map[string]string{})
	if len(slice) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(slice))
	}
}
