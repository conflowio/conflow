package function

import "strings"

// StringContains reports whether substr is within s.
//go:generate basil generate
func StringContains(s, substr string) bool {
	return strings.Contains(s, substr)
}
