package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTemp(t, `
vault_address: https://vault.example.com
mappings:
  - vault_path: secret/myapp
    env_file: .env
    mount: secret
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VaultAddress != "https://vault.example.com" {
		t.Errorf("expected vault address, got %q", cfg.VaultAddress)
	}
	if len(cfg.Mappings) != 1 {
		t.Fatalf("expected 1 mapping, got %d", len(cfg.Mappings))
	}
	if cfg.Mappings[0].VaultPath != "secret/myapp" {
		t.Errorf("unexpected vault_path: %q", cfg.Mappings[0].VaultPath)
	}
}

func TestLoad_EmptyPath(t *testing.T) {
	_, err := Load("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_NoMappings(t *testing.T) {
	path := writeTemp(t, `vault_address: https://vault.example.com\nmappings: []\n`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for empty mappings")
	}
}

func TestLoad_MissingVaultPath(t *testing.T) {
	path := writeTemp(t, `
mappings:
  - env_file: .env
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing vault_path")
	}
}
