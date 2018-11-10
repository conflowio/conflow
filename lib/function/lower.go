package function

import (
	"strings"
)

// Lower returns a copy of the string s with all Unicode letters mapped to their lower case.
//go:generate basil generate Lower
func Lower(s string) string {
	return strings.ToLower(s)
}
