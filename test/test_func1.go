package test

import (
	"strings"
)

//go:generate basil generate
func testFunc1(str string) string {
	return strings.ToUpper(str)
}
