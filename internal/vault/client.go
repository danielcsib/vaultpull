package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods.
type Client struct {
	logical *vaultapi.Logical
}

// Config holds configuration for connecting to Vault.
type Config struct {
	Address string
	Token   string
}

// NewClient creates a new Vault client from the given config.
// If config fields are empty, it falls back to environment variables.
func NewClient(cfg Config) (*Client, error) {
	vaultCfg := vaultapi.DefaultConfig()

	address := cfg.Address
	if address == "" {
		address = os.Getenv("VAULT_ADDR")
	}
	if address == "" {
		return nil, fmt.Errorf("vault address not set: provide --vault-addr or set VAULT_ADDR")
	}
	vaultCfg.Address = address

	client, err := vaultapi.NewClient(vaultCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	token := cfg.Token
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token not set: provide --vault-token or set VAULT_TOKEN")
	}
	client.SetToken(token)

	return &Client{logical: client.Logical()}, nil
}

// ReadSecrets reads key-value secrets from the given Vault path.
// Supports both KV v1 and KV v2 paths.
func (c *Client) ReadSecrets(path string) (map[string]string, error) {
	secret, err := c.logical.Read(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read path %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at path %q", path)
	}

	// KV v2 wraps data under a "data" key.
	data := secret.Data
	if nested, ok := data["data"]; ok {
		if nestedMap, ok := nested.(map[string]interface{}); ok {
			data = nestedMap
		}
	}

	result := make(map[string]string, len(data))
	for k, v := range data {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result, nil
}
