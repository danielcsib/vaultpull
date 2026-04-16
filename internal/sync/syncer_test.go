package sync_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/vaultpull/internal/config"
	"github.com/user/vaultpull/internal/sync"
	"github.com/user/vaultpull/internal/vault"
)

func newTestClient(t *testing.T, addr, token string) *vault.Client {
	t.Helper()
	c, err := vault.NewClient(addr, token)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c
}

func TestRun_EmptyMappings(t *testing.T) {
	client := newTestClient(t, "http://127.0.0.1:8200", "test-token")
	s := sync.New(client)
	cfg := &config.Config{Mappings: []config.Mapping{}}
	results := s.Run(cfg)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestRun_VaultUnreachable(t *testing.T) {
	client := newTestClient(t, "http://127.0.0.1:19999", "test-token")
	s := sync.New(client)

	tmpDir := t.TempDir()
	cfg := &config.Config{
		Mappings: []config.Mapping{
			{VaultPath: "secret/data/app", EnvFile: filepath.Join(tmpDir, ".env")},
		},
	}

	results := s.Run(cfg)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Err == nil {
		t.Error("expected error for unreachable vault, got nil")
	}
}

func TestRun_WritesEnvFile(t *testing.T) {
	// This test uses a pre-written secrets stub via the envwriter directly;
	// full integration requires a live Vault — skipped in unit suite.
	t.Skip("integration test: requires live Vault")

	client := newTestClient(t, os.Getenv("VAULT_ADDR"), os.Getenv("VAULT_TOKEN"))
	s := sync.New(client)
	tmpDir := t.TempDir()
	cfg := &config.Config{
		Mappings: []config.Mapping{
			{VaultPath: "secret/data/app", EnvFile: filepath.Join(tmpDir, ".env")},
		},
	}
	results := s.Run(cfg)
	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
}
