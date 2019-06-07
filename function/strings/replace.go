package strings

import "strings"

// Replace returns a copy of the string s with all
// non-overlapping instances of old replaced by new.
//go:generate basil generate
func Replace(s, old, new string) string {
	return strings.Replace(s, old, new, -1)
}
