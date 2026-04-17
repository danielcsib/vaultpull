package sync

import (
	"fmt"

	"github.com/yourorg/vaultpull/internal/cache"
	"github.com/yourorg/vaultpull/internal/config"
	"github.com/yourorg/vaultpull/internal/envwriter"
	"github.com/yourorg/vaultpull/internal/notify"
	"github.com/yourorg/vaultpull/internal/vault"
)

// VaultClient abstracts secret reading from Vault.
type VaultClient interface {
	ReadSecrets(path string) (map[string]string, error)
}

// Result holds the outcome of syncing a single mapping.
type Result struct {
	Path    string
	OutFile string
	Err     error
	Skipped bool
}

// Syncer orchestrates pulling secrets and writing env files.
type Syncer struct {
	client   VaultClient
	cache    *cache.Store
	notifier *notify.Notifier
}

// New constructs a Syncer.
func New(c VaultClient, s *cache.Store, n *notify.Notifier) *Syncer {
	return &Syncer{client: c, cache: s, notifier: n}
}

// Run executes all mappings defined in cfg and returns per-mapping results.
func (s *Syncer) Run(cfg *config.Config) []Result {
	results := make([]Result, 0, len(cfg.Mappings))
	for _, m := range cfg.Mappings {
		r := Result{Path: m.VaultPath, OutFile: m.EnvFile}
		secrets, err := s.client.ReadSecrets(m.VaultPath)
		if err != nil {
			r.Err = fmt.Errorf("read %s: %w", m.VaultPath, err)
			if s.notifier != nil {
				s.notifier.Error(m.VaultPath, r.Err.Error())
			}
			results = append(results, r)
			continue
		}
		if s.cache != nil {
			if changed, _ := s.cache.Changed(m.VaultPath, secrets); !changed {
				r.Skipped = true
				if s.notifier != nil {
					s.notifier.Warn(m.VaultPath, "no changes, skipping")
				}
				results = append(results, r)
				continue
			}
		}
		if err := envwriter.WriteEnvFile(m.EnvFile, secrets); err != nil {
			r.Err = fmt.Errorf("write %s: %w", m.EnvFile, err)
			if s.notifier != nil {
				s.notifier.Error(m.VaultPath, r.Err.Error())
			}
			results = append(results, r)
			continue
		}
		if s.cache != nil {
			_ = s.cache.Set(m.VaultPath, secrets)
		}
		if s.notifier != nil {
			s.notifier.Info(m.VaultPath, fmt.Sprintf("wrote %d keys to %s", len(secrets), m.EnvFile))
		}
		results = append(results, r)
	}
	return results
}

// ensure vault import is used indirectly via interface
var _ VaultClient = (*vault.Client)(nil)
