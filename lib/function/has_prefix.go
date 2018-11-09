package function

import (
	"strings"
)

// HasPrefix tests whether the string s begins with prefix.
//go:generate basil generate HasPrefix
func HasPrefix(s string, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}
