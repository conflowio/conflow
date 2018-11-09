package function

import "strings"

// Join concatenates the elements of a to create a single string. The separator string
// sep is placed between elements in the resulting string.
//go:generate basil generate Join
func Join(a []string, sep string) string {
	return strings.Join(a, sep)
}
