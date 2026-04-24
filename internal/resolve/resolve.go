// Package resolve provides path resolution utilities for Vault secret paths,
// supporting variable substitution and environment-based interpolation.
package resolve

import (
	"fmt"
	"os"
	"strings"
)

// Resolver substitutes variables in Vault secret paths using a provided
// environment map, falling back to process environment variables.
type Resolver struct {
	env map[string]string
}

// New returns a Resolver that uses the provided map for variable lookup.
// If env is nil, only the process environment is consulted.
func New(env map[string]string) *Resolver {
	if env == nil {
		env = make(map[string]string)
	}
	return &Resolver{env: env}
}

// Resolve replaces all occurrences of ${VAR} or $VAR in the given path
// with values from the resolver's env map, falling back to os.Getenv.
// Returns an error if a referenced variable is not found in either source.
func (r *Resolver) Resolve(path string) (string, error) {
	var err error
	result := os.Expand(path, func(key string) string {
		if val, ok := r.env[key]; ok {
			return val
		}
		if val, ok := os.LookupEnv(key); ok {
			return val
		}
		err = fmt.Errorf("resolve: variable %q not found", key)
		return ""
	})
	if err != nil {
		return "", err
	}
	return result, nil
}

// ResolveAll resolves a slice of paths, returning all resolved values or the
// first error encountered.
func (r *Resolver) ResolveAll(paths []string) ([]string, error) {
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		resolved, err := r.Resolve(p)
		if err != nil {
			return nil, err
		}
		out = append(out, resolved)
	}
	return out, nil
}

// HasVariables reports whether the given path contains any substitution
// expressions of the form ${VAR} or $VAR.
func HasVariables(path string) bool {
	return strings.Contains(path, "$")
}
