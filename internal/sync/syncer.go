package sync

import (
	"fmt"

	"github.com/user/vaultpull/internal/config"
	"github.com/user/vaultpull/internal/envwriter"
	"github.com/user/vaultpull/internal/vault"
)

// Result holds the outcome of a single mapping sync.
type Result struct {
	Mapping config.Mapping
	Written int
	Err     error
}

// Syncer orchestrates pulling secrets from Vault and writing .env files.
type Syncer struct {
	client *vault.Client
}

// New creates a Syncer with the provided Vault client.
func New(client *vault.Client) *Syncer {
	return &Syncer{client: client}
}

// Run iterates over all mappings in the config and syncs each one.
func (s *Syncer) Run(cfg *config.Config) []Result {
	results := make([]Result, 0, len(cfg.Mappings))
	for _, m := range cfg.Mappings {
		written, err := s.syncMapping(m)
		results = append(results, Result{Mapping: m, Written: written, Err: err})
	}
	return results
}

func (s *Syncer) syncMapping(m config.Mapping) (int, error) {
	secrets, err := s.client.ReadSecrets(m.VaultPath)
	if err != nil {
		return 0, fmt.Errorf("reading vault path %q: %w", m.VaultPath, err)
	}

	if err := envwriter.ValidateKeys(secrets); err != nil {
		return 0, fmt.Errorf("invalid keys at %q: %w", m.VaultPath, err)
	}

	if err := envwriter.WriteEnvFile(m.EnvFile, secrets); err != nil {
		return 0, fmt.Errorf("writing env file %q: %w", m.EnvFile, err)
	}

	return len(secrets), nil
}
