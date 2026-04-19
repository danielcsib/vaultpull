// Package schema validates that a set of secrets conforms to an expected shape.
package schema

import (
	"errors"
	"fmt"
	"strings"
)

// FieldRule describes expectations for a single secret key.
type FieldRule struct {
	Key      string
	Required bool
	AllowEmpty bool
}

// Schema holds the full set of rules for a secret path.
type Schema struct {
	rules []FieldRule
}

// Violation describes a single rule breach.
type Violation struct {
	Key     string
	Message string
}

func (v Violation) Error() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// New creates a Schema from the provided rules.
func New(rules []FieldRule) *Schema {
	return &Schema{rules: rules}
}

// Validate checks secrets against the schema and returns all violations.
func (s *Schema) Validate(secrets map[string]string) ([]Violation, error) {
	var violations []Violation

	for _, rule := range s.rules {
		val, exists := secrets[rule.Key]
		if rule.Required && !exists {
			violations = append(violations, Violation{Key: rule.Key, Message: "required key is missing"})
			continue
		}
		if exists && !rule.AllowEmpty && strings.TrimSpace(val) == "" {
			violations = append(violations, Violation{Key: rule.Key, Message: "value must not be empty"})
		}
	}

	if len(violations) > 0 {
		return violations, errors.New("schema validation failed")
	}
	return nil, nil
}
