// Package scope provides path-scoped secret filtering for vault mappings.
// It allows restricting which Vault secret paths are eligible for sync
// based on a configured root prefix.
package scope

import (
	"fmt"
	"strings"
)

// Scoper restricts secret paths to a configured root prefix.
type Scoper struct {
	root string
}

// New creates a Scoper with the given root prefix.
// The root must be non-empty and is normalised to have no trailing slash.
func New(root string) (*Scoper, error) {
	root = strings.TrimSpace(root)
	if root == "" {
		return nil, fmt.Errorf("scope: root prefix must not be empty")
	}
	root = strings.TrimRight(root, "/")
	return &Scoper{root: root}, nil
}

// Contains reports whether path falls within the scoper's root prefix.
func (s *Scoper) Contains(path string) bool {
	path = strings.TrimRight(path, "/")
	if path == s.root {
		return true
	}
	return strings.HasPrefix(path, s.root+"/")
}

// Filter returns only those paths from the input slice that are within scope.
func (s *Scoper) Filter(paths []string) []string {
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		if s.Contains(p) {
			out = append(out, p)
		}
	}
	return out
}

// Root returns the configured root prefix.
func (s *Scoper) Root() string {
	return s.root
}

// RelativePath strips the root prefix from path, returning the remainder.
// If path is not within scope, the original path is returned unchanged.
func (s *Scoper) RelativePath(path string) string {
	path = strings.TrimRight(path, "/")
	if path == s.root {
		return ""
	}
	if strings.HasPrefix(path, s.root+"/") {
		return strings.TrimPrefix(path, s.root+"/")
	}
	return path
}
