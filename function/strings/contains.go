package strings

import "strings"

// Contains reports whether substr is within s.
//go:generate basil generate
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
