package mask

import "strings"

// Sensitive is a list of key substrings that should be masked in output.
var defaultSensitivePatterns = []string{
	"password", "secret", "token", "key", "auth", "credential",
}

// Masker redacts sensitive keys from maps before display or logging.
type Masker struct {
	patterns []string
	placeholder string
}

// New returns a Masker with the default sensitive patterns.
func New() *Masker {
	return &Masker{
		patterns:    defaultSensitivePatterns,
		placeholder: "***REDACTED***",
	}
}

// NewWithPatterns returns a Masker using custom patterns.
func NewWithPatterns(patterns []string, placeholder string) *Masker {
	if placeholder == "" {
		placeholder = "***REDACTED***"
	}
	return &Masker{patterns: patterns, placeholder: placeholder}
}

// IsSensitive returns true if the key matches any sensitive pattern.
func (m *Masker) IsSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, p := range m.patterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// MaskMap returns a copy of the map with sensitive values replaced.
func (m *Masker) MaskMap(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if m.IsSensitive(k) {
			out[k] = m.placeholder
		} else {
			out[k] = v
		}
	}
	return out
}
