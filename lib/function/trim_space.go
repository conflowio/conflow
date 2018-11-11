package function

import (
	"strings"
)

// TrimSpace returns a slice of the string s, with all leading
// and trailing white space removed, as defined by Unicode.
//go:generate basil generate
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}
