package test

import (
	"strings"
)

//go:generate basil generate testFunc1
func testFunc1(str string) string {
	return strings.ToUpper(str)
}
