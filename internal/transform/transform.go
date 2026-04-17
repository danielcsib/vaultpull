package transform

import (
	"fmt"
	"strings"
)

// Rule defines a transformation to apply to a secret key or value.
type Rule struct {
	Prefix    string // prepend to key
	Uppercase bool   // force key to uppercase
	Strip     string // strip this prefix from the vault key
}

// Apply applies the transformation rule to a map of secrets,
// returning a new map with transformed keys.
func Apply(secrets map[string]string, rule Rule) (map[string]string, error) {
	if secrets == nil {
		return nil, fmt.Errorf("transform: secrets map must not be nil")
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		transformed, err := transformKey(k, rule)
		if err != nil {
			return nil, err
		}
		result[transformed] = v
	}
	return result, nil
}

// transformKey applies rule transformations to a single key.
func transformKey(key string, rule Rule) (string, error) {
	if key == "" {
		return "", fmt.Errorf("transform: empty key is not allowed")
	}

	if rule.Strip != "" {
		key = strings.TrimPrefix(key, rule.Strip)
		if key == "" {
			return "", fmt.Errorf("transform: key became empty after stripping prefix %q", rule.Strip)
		}
	}

	if rule.Uppercase {
		key = strings.ToUpper(key)
	}

	if rule.Prefix != "" {
		key = rule.Prefix + key
	}

	return key, nil
}
