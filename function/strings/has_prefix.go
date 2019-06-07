package strings

import (
	"strings"
)

// HasPrefix tests whether the string s begins with prefix.
//go:generate basil generate
func HasPrefix(s string, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}
