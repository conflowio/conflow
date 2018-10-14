package test

import "github.com/onsi/ginkgo/extensions/table"

// TableEntry creates an custom entry for table driven tests where the input is the description
func TableEntry(input string, parameters ...interface{}) table.TableEntry {
	return table.Entry(input, append([]interface{}{input}, parameters...)...)
}
