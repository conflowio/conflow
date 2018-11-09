package function

import (
	"strings"
)

// HasSuffix tests whether the string s ends with suffix.
//go:generate basil generate HasSuffix
func HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}
