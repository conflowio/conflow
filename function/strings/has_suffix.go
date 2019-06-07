package strings

import (
	"strings"
)

// HasSuffix tests whether the string s ends with suffix.
//go:generate basil generate
func HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}
