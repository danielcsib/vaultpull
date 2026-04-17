package sync

import (
	"fmt"

	"github.com/yourusername/vaultpull/internal/cache"
	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/envwriter"
	"github.com/yourusername/vaultpull/internal/vault"
)

// Result holds the outcome of syncing a single mapping.
type Result struct {
	VaultPath string
	EnvFile   string
	Skipped   bool
	Err       error
}

// Syncer orchestrates fetching secrets and writing env files.
type Syncer struct {
	client *vault.Client
	cache  *cache.Store
}

// New creates a Syncer with an optional cache store (may be nil).
func New(client *vault.Client, store *cache.Store) *Syncer {
	return &Syncer{client: client, cache: store}
}

// Run processes all mappings defined in cfg and returns per-mapping results.
func (s *Syncer) Run(cfg *config.Config) []Result {
	results := make([]Result, 0, len(cfg.Mappings))
	for _, m := range cfg.Mappings {
		results = append(results, s.sync(m))
	}
	return results
}

func (s *Syncer) sync(m config.Mapping) Result {
	r := Result{VaultPath: m.VaultPath, EnvFile: m.EnvFile}

	secrets, err := s.client.ReadSecrets(m.VaultPath)
	if err != nil {
		r.Err = fmt.Errorf("read %s: %w", m.VaultPath, err)
		return r
	}

	if s.cache != nil {
		changed, err := s.cache.Changed(m.VaultPath, secrets)
		if err == nil && !changed {
			r.Skipped = true
			return r
		}
	}

	if err := envwriter.WriteEnvFile(m.EnvFile, secrets); err != nil {
		r.Err = fmt.Errorf("write %s: %w", m.EnvFile, err)
		return r
	}

	if s.cache != nil {
		_ = s.cache.Set(m.VaultPath, secrets)
	}
	return r
}
