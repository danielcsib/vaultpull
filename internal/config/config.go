package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

// SecretMapping defines a single Vault path to env file mapping.
type SecretMapping struct {
	VaultPath string `yaml:"vault_path"`
	EnvFile   string `yaml:"env_file"`
	Mount     string `yaml:"mount"`
}

// Config represents the top-level vaultpull configuration.
type Config struct {
	VaultAddress string          `yaml:"vault_address"`
	Mappings     []SecretMapping `yaml:"mappings"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("config path must not be empty")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening config file: %w", err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if len(c.Mappings) == 0 {
		return errors.New("config must define at least one mapping")
	}
	for i, m := range c.Mappings {
		if m.VaultPath == "" {
			return fmt.Errorf("mapping[%d]: vault_path must not be empty", i)
		}
		if m.EnvFile == "" {
			return fmt.Errorf("mapping[%d]: env_file must not be empty", i)
		}
	}
	return nil
}
