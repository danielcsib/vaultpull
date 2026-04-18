// Package merge combines existing env file contents with newly pulled secrets,
// preserving local-only keys and updating changed values.
package merge

import (
	"bufio"
	"os"
	"strings"
)

// Result holds the outcome of a merge operation.
type Result struct {
	Final    map[string]string
	Added    []string
	Updated  []string
	Preserved []string
}

// FromFile reads an existing .env file and merges it with incoming secrets.
// Keys present only in the file are preserved. Incoming secrets take precedence
// for overlapping keys.
func FromFile(path string, incoming map[string]string) (Result, error) {
	existing, err := parseEnvFile(path)
	if err != nil && !os.IsNotExist(err) {
		return Result{}, err
	}

	final := make(map[string]string, len(existing)+len(incoming))
	result := Result{Final: final}

	for k, v := range existing {
		final[k] = v
	}

	for k, v := range incoming {
		if _, exists := existing[k]; !exists {
			result.Added = append(result.Added, k)
		} else if existing[k] != v {
			result.Updated = append(result.Updated, k)
		}
		final[k] = v
	}

	for k := range existing {
		if _, inIncoming := incoming[k]; !inIncoming {
			result.Preserved = append(result.Preserved, k)
		}
	}

	return result, nil
}

func parseEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	out := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		out[key] = val
	}
	return out, scanner.Err()
}
