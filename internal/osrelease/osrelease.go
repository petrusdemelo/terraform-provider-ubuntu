// Package osrelease reads the os-release(5) file.
package osrelease

import (
	"bufio"
	"strconv"
	"strings"
)

const Command = "cat /etc/os-release"

// Parse returns the KEY=VALUE assignments in an os-release file. Values may
// be unquoted, single-quoted or double-quoted; blank lines and comments are
// skipped, as are lines that are not assignments.
func Parse(content string) map[string]string {
	fields := make(map[string]string)

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok || key == "" {
			continue
		}

		fields[key] = unquote(value)
	}

	return fields
}

func unquote(value string) string {
	if len(value) >= 2 && value[0] == '\'' && value[len(value)-1] == '\'' {
		return value[1 : len(value)-1]
	}

	if unquoted, err := strconv.Unquote(value); err == nil && strings.HasPrefix(value, `"`) {
		return unquoted
	}

	return value
}
