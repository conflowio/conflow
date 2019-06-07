package strings

import (
	"strings"
)

// Lower returns a copy of the string s with all Unicode letters mapped to their lower case.
//go:generate basil generate
func Lower(s string) string {
	return strings.ToLower(s)
}
