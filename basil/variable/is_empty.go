package variable

import (
	"fmt"
	"time"

	"github.com/opsidian/basil/basil"
)

// IsEmpty returns true if the given value has an empty value
func IsEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case []interface{}:
		return len(v) == 0
	case bool:
		return v == false
	case float64:
		return v == 0.0
	case basil.ID:
		return string(v) == ""
	case int64:
		return v == int64(0)
	case map[string]interface{}:
		return len(v) == 0
	case string:
		return v == ""
	case []string:
		return len(v) == 0
	case time.Time:
		return v.IsZero()
	case time.Duration:
		return v == 0
	case Union:
		return IsEmpty(v.Value())
	default:
		panic(fmt.Sprintf("unexpected type: %T", value))
	}
}
