package function

import (
	"fmt"
	"unicode/utf8"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/variable"
)

// Len returns with the length of the variable
// For strings it means the count of UTF-8 characters
// For arrays and maps it means the number of items/entries
//go:generate basil generate
func Len(value *variable.WithLength) int64 {
	switch v := value.Value().(type) {
	case string:
		return int64(utf8.RuneCountInString(v))
	case basil.ID:
		return int64(utf8.RuneCountInString(string(v)))
	case []interface{}:
		return int64(len(v))
	case []string:
		return int64(len(v))
	case map[string]interface{}:
		return int64(len(v))
	default:
		panic(fmt.Sprintf("unexpected type: %T", v))
	}
}
