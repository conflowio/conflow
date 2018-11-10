package function

import (
	"strings"
)

// Trim returns a slice of the string s, with all leading
// and trailing white space removed, as defined by Unicode.
//go:generate basil generate Trim
func Trim(s string) string {
	return strings.TrimSpace(s)
}
