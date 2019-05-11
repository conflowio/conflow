package function

import "github.com/opsidian/basil/basil/variable"

// ArrayContains returns true if the array contains the given element
//go:generate basil generate
func ArrayContains(arr []interface{}, elem interface{}) bool {
	switch elem.(type) {
	case []interface{}, map[string]interface{}:
		for _, item := range arr {
			if variable.Equals(elem, item) {
				return true
			}
		}
	default:
		for _, item := range arr {
			if item == elem {
				return true
			}
		}
	}

	return false
}
