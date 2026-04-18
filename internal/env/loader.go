package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Load reads an .env file and returns a map of key-value pairs.
// Lines starting with '#' and empty lines are ignored.
// Values may be optionally quoted with single or double quotes.
func Load(path string) (map[string]string, error) {
	if path == "" {
		return nil, fmt.Errorf("env: path must not be empty")
	}

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, fmt.Errorf("env: open %s: %w", path, err)
	}
	defer f.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(f)
	lineNo := 0

	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 1 {
			return nil, fmt.Errorf("env: %s line %d: invalid format", path, lineNo)
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		val = unquote(val)
		result[key] = val
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("env: scan %s: %w", path, err)
	}
	return result, nil
}

func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
