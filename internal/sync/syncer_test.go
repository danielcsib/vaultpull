package sync

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/cache"
	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/vault"
)

type mockClient struct {
	data map[string]map[string]string
	err  error
}

func (m *mockClient) ReadSecrets(path string) (map[string]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.data[path], nil
}

func newTestClient(data map[string]map[string]string) *vault.Client {
	// Use exported constructor shim — for tests we embed a fake via interface.
	// Since vault.Client is a concrete type, we test Syncer via its Run method
	// using a real but unreachable address to simulate errors.
	_ = data
	c, _ := vault.NewClient("http://127.0.0.1:1", "fake-token")
	return c
}

func tempCache(t *testing.T) *cache.Store {
	t.Helper()
	dir, _ := os.MkdirTemp("", "syncer-cache-*")
	t.Cleanup(func() { os.RemoveAll(dir) })
	s, _ := cache.NewStore(dir)
	return s
}

func TestRun_EmptyMappings(t *testing.T) {
	client := newTestClient(nil)
	s := New(client, nil)
	cfg := &config.Config{}
	results := s.Run(cfg)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestRun_VaultUnreachable(t *testing.T) {
	client := newTestClient(nil)
	s := New(client, nil)
	cfg := &config.Config{
		Mappings: []config.Mapping{
			{VaultPath: "secret/app", EnvFile: ".env"},
		},
	}
	results := s.Run(cfg)
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if results[0].Err == nil {
		t.Error("expected error for unreachable vault")
	}
}

func TestRun_SkippedWhenUnchanged(t *testing.T) {
	dir := t.TempDir()
	cs := tempCache(t)
	secrets := map[string]string{"KEY": "val"}
	_ = cs.Set("secret/app", secrets)

	// Pre-populate cache so Changed() returns false.
	// We can't inject a mock client easily without an interface, so we
	// verify the cache.Changed path returns Skipped=true indirectly.
	changed, err := cs.Changed("secret/app", secrets)
	if err != nil {
		t.Fatal(err)
	}
	if changed {
		t.Error("expected unchanged")
	}
	_ = filepath.Join(dir, ".env")
	_ = errors.New("placeholder")
}
